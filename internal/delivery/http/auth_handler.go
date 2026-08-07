package http

import (
	"net/http"

	"backend_go/internal/models"
	"backend_go/internal/usecase"

	"backend_go/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase
}

func NewAuthHandler(api *gin.RouterGroup, protectedAPI *gin.RouterGroup, authUC usecase.AuthUseCase) {
	handler := &AuthHandler{
		authUC: authUC,
	}

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", handler.Login)
		authGroup.POST("/register", handler.Register)
		authGroup.POST("/forgot-password", handler.ForgotPassword)
		authGroup.POST("/reset-password", handler.ResetPassword)
	}

	protectedAuthGroup := protectedAPI.Group("/auth")
	{
		protectedAuthGroup.GET("/profile", handler.Profile)
		protectedAuthGroup.POST("/logout", handler.Logout)
		protectedAuthGroup.POST("/change-password", handler.ChangePassword)
	}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	token, user, err := h.authUC.Login(req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login berhasil", gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.authUC.Register(&user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", nil)
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	user, err := h.authUC.Profile(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", user)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, _ := c.Get("token")
	
	if err := h.authUC.Logout(token.(string)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Successfully logged out", nil)
}

type forgotPasswordReq struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req forgotPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.authUC.ForgotPassword(req.Email); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reset link sent to your email", nil)
}

type resetPasswordReq struct {
	Email    string `json:"email" binding:"required,email"`
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req resetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.authUC.ResetPassword(req.Email, req.Token, req.Password); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password reset successfully", nil)
}

type changePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req changePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	userID, _ := c.Get("userID")

	if err := h.authUC.ChangePassword(userID.(string), req.OldPassword, req.NewPassword); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password changed successfully", nil)
}
