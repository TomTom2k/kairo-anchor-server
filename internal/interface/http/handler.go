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
	resetPassword  *auth.ResetPasswordUseCase
}


func NewHandler(
	r *auth.RegisterUseCase,
	l *auth.LoginUseCase,
	gp *auth.GetProfileUseCase,
	a *auth.ActivateAccountUseCase,
	fp *auth.ForgotPasswordUseCase,
	cp *auth.ChangePasswordUseCase,
	rp *auth.ResetPasswordUseCase,
) *Handler {
	return &Handler{
		register:       r,
		login:          l,
		getProfile:     gp,
		activate:       a,
		forgotPassword: fp,
		changePassword: cp,
		resetPassword:  rp,
	}
}


// Register godoc
// @Summary Register a new user
// @Description Create a new user account and send activation email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register Request"
// @Success 201 {object} APIResponse
// @Failure 400 {object} APIErrorResponse
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	err := h.register.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		SendError(c, http.StatusBadRequest, "REGISTRATION_FAILED", err.Error())
		return
	}

	SendSuccess(c, http.StatusCreated, nil, "Registration successful, please check your email to activate your account")
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	result, err := h.login.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid email or password")
		return
	}

	SendSuccess(c, http.StatusOK, gin.H{
		"token": result.Token,
		"user": ProfileResponse{
			ID:        result.User.ID,
			Email:     result.User.Email,
			IsActive:  result.User.IsActive,
			CreatedAt: result.User.CreatedAt,
			UpdatedAt: result.User.UpdatedAt,
		},
	}, "Login successful")
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get authenticated user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} APIResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Router /auth/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized")
		return
	}

	user, err := h.getProfile.Execute(c.Request.Context(), userID)
	if err != nil {
		SendInternalError(c, err)
		return
	}

	if user == nil {
		SendError(c, http.StatusNotFound, ErrCodeNotFound, "User not found")
		return
	}

	SendSuccess(c, http.StatusOK, ProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, "")
}

// ActivateAccount godoc
// @Summary Activate user account
// @Description Activate user account using activation token from email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ActivateAccountRequest true "Activate Account Request"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIErrorResponse
// @Router /auth/activate [post]
func (h *Handler) ActivateAccount(c *gin.Context) {
	var req ActivateAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	err := h.activate.Execute(c.Request.Context(), req.Token)
	if err != nil {
		SendError(c, http.StatusBadRequest, "ACTIVATION_FAILED", err.Error())
		return
	}

	SendSuccess(c, http.StatusOK, nil, "Account activated successfully, you can now login")
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset email to user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Forgot Password Request"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIErrorResponse
// @Router /auth/forgot-password [post]
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	err := h.forgotPassword.Execute(c.Request.Context(), req.Email)
	if err != nil {
		SendError(c, http.StatusBadRequest, "PASSWORD_RESET_FAILED", err.Error())
		return
	}

	SendSuccess(c, http.StatusOK, nil, "Password reset email sent, please check your email")
}

// ChangePassword godoc
// @Summary Change password
// @Description Change password using reset token from email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Change Password Request"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIErrorResponse
// @Router /auth/change-password [post]
func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	err := h.changePassword.Execute(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		SendError(c, http.StatusBadRequest, "PASSWORD_CHANGE_FAILED", err.Error())
		return
	}

	SendSuccess(c, http.StatusOK, nil, "Password changed successfully, you can now login with your new password")
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password while logged in (requires old password)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Router /auth/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized")
		return
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	err = h.resetPassword.Execute(c.Request.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		SendError(c, http.StatusBadRequest, "PASSWORD_RESET_FAILED", err.Error())
		return
	}

	SendSuccess(c, http.StatusOK, nil, "Password reset successfully")
}