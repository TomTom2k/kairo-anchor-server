package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tomtom2k/kairo-anchor-server/internal/usecase/auth"
)

type Handler struct {
	register *auth.RegisterUseCase
}

func NewHandler(r *auth.RegisterUseCase) *Handler {
	return &Handler{r}
}

func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.register.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}
		