package http

import (
	"net/http"

	"backend_go/internal/usecase"
	"backend_go/internal/models"
	"backend_go/internal/middleware"
	"backend_go/pkg/utils"

	"github.com/gin-gonic/gin"
)

type MasterHandler struct {
	masterUC usecase.MasterUseCase
}

func NewMasterHandler(api *gin.RouterGroup, masterUC usecase.MasterUseCase) {
	handler := &MasterHandler{
		masterUC: masterUC,
	}

	adminOnly := middleware.RoleMiddleware(string(models.RoleAdmin), string(models.RoleOwner))
	api.GET("/patient-categories", handler.GetPatientCategories)
	api.POST("/patient-categories", adminOnly, handler.CreatePatientCategory)
	api.PUT("/patient-categories/:id", adminOnly, handler.UpdatePatientCategory)
	api.DELETE("/patient-categories/:id", adminOnly, handler.DeletePatientCategory)
	
	api.GET("/genders", handler.GetGenders)
}

func (h *MasterHandler) GetPatientCategories(c *gin.Context) {
	categories, err := h.masterUC.GetPatientCategories()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", categories)
}

func (h *MasterHandler) CreatePatientCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	category := models.PatientCategory{
		Name: req.Name,
	}

	if err := h.masterUC.CreatePatientCategory(&category); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Category created successfully", category)
}

func (h *MasterHandler) UpdatePatientCategory(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	category := models.PatientCategory{
		Name: req.Name,
	}

	if err := h.masterUC.UpdatePatientCategory(id, &category); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category updated successfully", category)
}

func (h *MasterHandler) DeletePatientCategory(c *gin.Context) {
	id := c.Param("id")

	if err := h.masterUC.DeletePatientCategory(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category deleted successfully", nil)
}

func (h *MasterHandler) GetGenders(c *gin.Context) {
	genders, err := h.masterUC.GetGenders()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", genders)
}
