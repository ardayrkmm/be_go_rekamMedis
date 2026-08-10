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

type PatientHandler struct {
	patientUC usecase.PatientUseCase
}

func NewPatientHandler(api *gin.RouterGroup, patientUC usecase.PatientUseCase) {
	handler := &PatientHandler{
		patientUC: patientUC,
	}

	group := api.Group("/patients")
	adminOnly := middleware.RoleMiddleware(string(models.RoleAdmin), string(models.RoleOwner))
	{
		group.GET("", handler.Fetch)
		group.GET("/:id", handler.GetByID)
		group.POST("", adminOnly, handler.Store)
		group.PUT("/:id", adminOnly, handler.Update)
		group.DELETE("/:id", adminOnly, handler.Delete)
		group.POST("/:id/restore", adminOnly, handler.Restore)
	}
}

func (h *PatientHandler) Fetch(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))
	search := c.Query("search")

	offset := (page - 1) * perPage

	patients, total, err := h.patientUC.Fetch(offset, perPage, search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponsePaginated(c, http.StatusOK, "Success", patients, page, perPage, total)
}

func (h *PatientHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	patient, err := h.patientUC.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Patient not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", patient)
}

func (h *PatientHandler) Store(c *gin.Context) {
	var patient models.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.patientUC.Store(&patient); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success", patient)
}

func (h *PatientHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var patient models.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.patientUC.Update(id, &patient); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Patient updated successfully", nil)
}

func (h *PatientHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.patientUC.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Patient deleted successfully", nil)
}

func (h *PatientHandler) Restore(c *gin.Context) {
	id := c.Param("id")
	if err := h.patientUC.Restore(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Patient restored successfully", nil)
}
