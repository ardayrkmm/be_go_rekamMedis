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

type ServiceMasterHandler struct {
	serviceUC usecase.ServiceMasterUseCase
}

func NewServiceMasterHandler(api *gin.RouterGroup, serviceUC usecase.ServiceMasterUseCase) {
	handler := &ServiceMasterHandler{
		serviceUC: serviceUC,
	}

	group := api.Group("/service-masters")
	adminOnly := middleware.RoleMiddleware(string(models.RoleAdmin), string(models.RoleOwner))
	{
		group.GET("", handler.Fetch)
		group.GET("/:id", handler.GetByID)
		group.POST("", adminOnly, handler.Store)
		group.PUT("/:id", adminOnly, handler.Update)
		group.DELETE("/:id", adminOnly, handler.Delete)
	}
}

func (h *ServiceMasterHandler) Fetch(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))
	search := c.Query("search")

	offset := (page - 1) * perPage

	services, total, err := h.serviceUC.Fetch(offset, perPage, search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponsePaginated(c, http.StatusOK, "Success", services, page, perPage, total)
}

func (h *ServiceMasterHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	service, err := h.serviceUC.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Service Master not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", service)
}

func (h *ServiceMasterHandler) Store(c *gin.Context) {
	var service models.ServiceMaster
	if err := c.ShouldBindJSON(&service); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.serviceUC.Store(&service); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success", service)
}

func (h *ServiceMasterHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var service models.ServiceMaster
	if err := c.ShouldBindJSON(&service); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.serviceUC.Update(id, &service); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Service Master updated successfully", nil)
}

func (h *ServiceMasterHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.serviceUC.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Service Master deleted successfully", nil)
}
