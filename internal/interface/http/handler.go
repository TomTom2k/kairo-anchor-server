package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tomtom2k/kairo-anchor-server/internal/usecase/auth"
)

type Handler struct {
	register       *auth.RegisterUseCase
	login          *auth.LoginUseCase
	getProfile     *auth.GetProfileUseCase
	activate       *auth.ActivateAccountUseCase
	forgotPassword *auth.ForgotPasswordUseCase
	changePassword *auth.ChangePasswordUseCase
}

func NewHandler(
	r *auth.RegisterUseCase,
	l *auth.LoginUseCase,
	gp *auth.GetProfileUseCase,
	a *auth.ActivateAccountUseCase,
	fp *auth.ForgotPasswordUseCase,
	cp *auth.ChangePasswordUseCase,
) *Handler {
	return &Handler{
		register:       r,
		login:          l,
		getProfile:     gp,
		activate:       a,
		forgotPassword: fp,
		changePassword: cp,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account and send activation email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register Request"
// @Success 201 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.register.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, MessageResponse{
		Message: "Registration successful, please check your email to activate your account",
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	result, err := h.login.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: result.Token,
		User: ProfileResponse{
			ID:        result.User.ID,
			Email:     result.User.Email,
			IsActive:  result.User.IsActive,
			CreatedAt: result.User.CreatedAt,
			UpdatedAt: result.User.UpdatedAt,
		},
		Message: "Login successful",
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get authenticated user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ProfileResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /auth/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	user, err := h.getProfile.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	c.JSON(http.StatusOK, ProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// ActivateAccount godoc
// @Summary Activate user account
// @Description Activate user account using activation token from email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ActivateAccountRequest true "Activate Account Request"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/activate [post]
func (h *Handler) ActivateAccount(c *gin.Context) {
	var req ActivateAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.activate.Execute(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Account activated successfully, you can now login",
	})
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset email to user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Forgot Password Request"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/forgot-password [post]
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.forgotPassword.Execute(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Password reset email sent, please check your email",
	})
}

// ChangePassword godoc
// @Summary Change password
// @Description Change password using reset token from email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Change Password Request"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/change-password [post]
func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.changePassword.Execute(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Password changed successfully, you can now login with your new password",
	})
}