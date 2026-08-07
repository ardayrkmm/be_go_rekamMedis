package http

import (
	"net/http"
	"strconv"

	"backend_go/internal/models"
	"backend_go/internal/usecase"
	"backend_go/internal/middleware"

	"backend_go/pkg/utils"
	"github.com/gin-gonic/gin"
	"cloud.google.com/go/firestore"
)

type UserHandler struct {
	userUC usecase.UserUseCase
}

func NewUserHandler(api *gin.RouterGroup, userUC usecase.UserUseCase, db *firestore.Client) {
	handler := &UserHandler{
		userUC: userUC,
	}

	group := api.Group("/users")
	group.Use(middleware.AuthMiddleware(db), middleware.RoleMiddleware(string(models.RoleAdmin), string(models.RoleOwner)))
	{
		group.GET("", handler.Fetch)
		group.GET("/:id", handler.GetByID)
		group.POST("", handler.Store)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", handler.Delete)
		group.POST("/:id/restore", handler.Restore)
	}
}

func (h *UserHandler) Fetch(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))

	offset := (page - 1) * perPage

	users, total, err := h.userUC.Fetch(offset, perPage)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponsePaginated(c, http.StatusOK, "Success", users, page, perPage, total)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userUC.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", user)
}

func (h *UserHandler) Store(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.userUC.Store(&user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success", user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.userUC.Update(id, &user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User updated successfully", nil)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.userUC.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}

func (h *UserHandler) Restore(c *gin.Context) {
	id := c.Param("id")
	if err := h.userUC.Restore(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User restored successfully", nil)
}
