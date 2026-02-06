package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/tomtom2k/kairo-anchor-server/internal/infrastructure/postgres"
	"github.com/tomtom2k/kairo-anchor-server/internal/interface/http"
	"github.com/tomtom2k/kairo-anchor-server/internal/usecase/auth"
	"github.com/tomtom2k/kairo-anchor-server/pkg/crypto"
)

func main() {
	db, _ := pgx.Connect(context.Background(), "postgres://...")

	userRepo := postgres.NewUserRepository(db)
	hasher := crypto.BcryptHasher{}

	registerUC := auth.NewRegisterUseCase(userRepo, hasher)
	handler := http.NewHandler(registerUC)

	r := gin.Default()
	r.POST("/register", handler.Register)
	r.Run()
}
