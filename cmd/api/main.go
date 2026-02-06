// @title Kairo Anchor Authentication API
// @version 1.0
// @description Authentication API with JWT tokens for Kairo Anchor application
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

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/tomtom2k/kairo-anchor-server/internal/config"
	"github.com/tomtom2k/kairo-anchor-server/internal/infrastructure/postgres"
	"github.com/tomtom2k/kairo-anchor-server/internal/interface/http"
	"github.com/tomtom2k/kairo-anchor-server/internal/usecase/auth"
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

	// Initialize services
	userRepo := postgres.NewUserRepository(db)
	hasher := &crypto.BcryptHasher{}
	tokenService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	emailService := email.NewMockEmailService(cfg.App.BaseURL)

	// Initialize use cases
	registerUC := auth.NewRegisterUseCase(userRepo, hasher, emailService)
	loginUC := auth.NewLoginUseCase(userRepo, hasher, tokenService)
	getProfileUC := auth.NewGetProfileUseCase(userRepo)
	activateUC := auth.NewActivateAccountUseCase(userRepo)
	forgotPasswordUC := auth.NewForgotPasswordUseCase(userRepo, emailService)
	changePasswordUC := auth.NewChangePasswordUseCase(userRepo, hasher)

	// Initialize HTTP handler and middleware
	handler := http.NewHandler(registerUC, loginUC, getProfileUC, activateUC, forgotPasswordUC, changePasswordUC)
	authMiddleware := http.NewAuthMiddleware(tokenService)

	// Setup Gin router
	r := gin.Default()

	// Add recovery middleware to catch panics
	r.Use(http.Recovery())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			// Public routes
			authGroup.POST("/register", handler.Register)
			authGroup.POST("/login", handler.Login)
			authGroup.POST("/activate", handler.ActivateAccount)
			authGroup.POST("/forgot-password", handler.ForgotPassword)
			authGroup.POST("/change-password", handler.ChangePassword)

			// Protected routes
			authGroup.GET("/profile", authMiddleware.RequireAuth(), handler.GetProfile)
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
