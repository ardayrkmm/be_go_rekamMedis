package http

import (
	"net/http"
	"strconv"

	"backend_go/internal/models"
	"backend_go/internal/usecase"

	"backend_go/pkg/utils"
	"github.com/gin-gonic/gin"
)

type MedicalRecordHandler struct {
	recordUC usecase.MedicalRecordUseCase
}

func NewMedicalRecordHandler(api *gin.RouterGroup, recordUC usecase.MedicalRecordUseCase) {
	handler := &MedicalRecordHandler{
		recordUC: recordUC,
	}

	group := api.Group("/medical-records")
	{
		group.GET("", handler.Fetch)
		group.GET("/:id", handler.GetByID)
		group.POST("", handler.Store)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", handler.Delete)
	}

	// Histori Rekam Medis Pasien
	api.GET("/patients/:id/medical-records", handler.GetHistory)
}

func (h *MedicalRecordHandler) Fetch(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))

	offset := (page - 1) * perPage

	records, total, err := h.recordUC.Fetch(offset, perPage)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponsePaginated(c, http.StatusOK, "Success", records, page, perPage, total)
}

func (h *MedicalRecordHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	record, err := h.recordUC.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Medical Record not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", record)
}

func (h *MedicalRecordHandler) GetHistory(c *gin.Context) {
	patientID := c.Param("id")
	records, err := h.recordUC.GetByPatientID(patientID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", records)
}

func (h *MedicalRecordHandler) Store(c *gin.Context) {
	var record models.MedicalRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.recordUC.Store(&record); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success", record)
}

func (h *MedicalRecordHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var record models.MedicalRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.recordUC.Update(id, &record); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Medical Record updated successfully", nil)
}

func (h *MedicalRecordHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.recordUC.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Medical Record deleted successfully", nil)
}
