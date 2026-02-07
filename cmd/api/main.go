// @title Kairo Anchor API
// @version 1.0
// @description API with authentication and Project management for Kairo Anchor application
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@kairo-anchor.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/tomtom2k/kairo-anchor-server/internal/config"
	"github.com/tomtom2k/kairo-anchor-server/internal/infrastructure/postgres"
	"github.com/tomtom2k/kairo-anchor-server/internal/interface/http"
	"github.com/tomtom2k/kairo-anchor-server/internal/usecase/auth"
	projectUC "github.com/tomtom2k/kairo-anchor-server/internal/usecase/project"
	"github.com/tomtom2k/kairo-anchor-server/pkg/crypto"
	"github.com/tomtom2k/kairo-anchor-server/pkg/email"
	"github.com/tomtom2k/kairo-anchor-server/pkg/jwt"

	_ "github.com/tomtom2k/kairo-anchor-server/docs" // Import generated docs
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Setup database connection
	db, err := sql.Open("pgx", cfg.DatabaseURL())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.PingContext(context.Background()); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("✓ Database connected successfully")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	projectRepo := postgres.NewProjectRepository(db)

	// Initialize services
	hasher := &crypto.BcryptHasher{}
	tokenService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	emailService := email.NewMockEmailService(cfg.App.BaseURL)

	// Initialize auth use cases
	registerUC := auth.NewRegisterUseCase(userRepo, hasher, emailService)
	loginUC := auth.NewLoginUseCase(userRepo, hasher, tokenService)
	getProfileUC := auth.NewGetProfileUseCase(userRepo)
	activateUC := auth.NewActivateAccountUseCase(userRepo)
	forgotPasswordUC := auth.NewForgotPasswordUseCase(userRepo, emailService)
	changePasswordUC := auth.NewChangePasswordUseCase(userRepo, hasher)
	resetPasswordUC := auth.NewResetPasswordUseCase(userRepo, hasher)

	// Initialize project use cases
	createProjectUC := projectUC.NewCreateProjectUseCase(projectRepo)
	updateProjectUC := projectUC.NewUpdateProjectUseCase(projectRepo)
	deleteProjectUC := projectUC.NewDeleteProjectUseCase(projectRepo)
	getProjectUC := projectUC.NewGetProjectUseCase(projectRepo)
	listProjectsUC := projectUC.NewListProjectsUseCase(projectRepo)

	// Initialize HTTP handlers
	authHandler := http.NewHandler(registerUC, loginUC, getProfileUC, activateUC, forgotPasswordUC, changePasswordUC, resetPasswordUC)
	projectHandler := http.NewProjectHandler(createProjectUC, updateProjectUC, deleteProjectUC, getProjectUC, listProjectsUC)

	authMiddleware := http.NewAuthMiddleware(tokenService)

	// Setup Gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Add recovery middleware to catch panics
	r.Use(http.Recovery())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api")
	{
		// Auth routes
		authGroup := api.Group("/auth")
		{
			// Public routes
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/activate", authHandler.ActivateAccount)
			authGroup.POST("/forgot-password", authHandler.ForgotPassword)
			authGroup.POST("/change-password", authHandler.ChangePassword)

			// Protected routes
			authGroup.GET("/profile", authMiddleware.RequireAuth(), authHandler.GetProfile)
			authGroup.POST("/reset-password", authMiddleware.RequireAuth(), authHandler.ResetPassword)
		}

		// Project routes (all protected)
		projectGroup := api.Group("/projects", authMiddleware.RequireAuth())
		{
			projectGroup.POST("", projectHandler.CreateProject)
			projectGroup.GET("", projectHandler.ListProjects)
			projectGroup.GET("/:id", projectHandler.GetProject)
			projectGroup.PUT("/:id", projectHandler.UpdateProject)
			projectGroup.DELETE("/:id", projectHandler.DeleteProject)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		http.SendSuccess(c, 200, gin.H{"status": "ok"}, "Service is healthy")
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("🚀 Server starting on http://localhost%s", addr)
	log.Printf("📚 Swagger UI available at http://localhost%s/swagger/index.html", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
