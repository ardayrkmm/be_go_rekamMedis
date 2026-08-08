package http

import (
	"net/http"
	"strconv"

	"backend_go/internal/models"
	"backend_go/internal/usecase"
	"backend_go/internal/middleware"

	"backend_go/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ServiceCategoryHandler struct {
	categoryUC usecase.ServiceCategoryUseCase
}

func NewServiceCategoryHandler(api *gin.RouterGroup, categoryUC usecase.ServiceCategoryUseCase) {
	handler := &ServiceCategoryHandler{
		categoryUC: categoryUC,
	}

	group := api.Group("/service-categories")
	adminOnly := middleware.RoleMiddleware(string(models.RoleAdmin), string(models.RoleOwner))
	{
		group.GET("", handler.Fetch)
		group.GET("/:id", handler.GetByID)
		group.POST("", adminOnly, handler.Store)
		group.PUT("/:id", adminOnly, handler.Update)
		group.DELETE("/:id", adminOnly, handler.Delete)
	}
}

func (h *ServiceCategoryHandler) Fetch(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	offset := (page - 1) * perPage

	categories, total, err := h.categoryUC.Fetch(offset, perPage)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponsePaginated(c, http.StatusOK, "Success", categories, page, perPage, total)
}

func (h *ServiceCategoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	category, err := h.categoryUC.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Service Category not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", category)
}

func (h *ServiceCategoryHandler) Store(c *gin.Context) {
	var category models.ServiceCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.categoryUC.Store(&category); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success", category)
}

func (h *ServiceCategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var category models.ServiceCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.categoryUC.Update(id, &category); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Service Category updated successfully", nil)
}

func (h *ServiceCategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.categoryUC.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Service Category deleted successfully", nil)
}
