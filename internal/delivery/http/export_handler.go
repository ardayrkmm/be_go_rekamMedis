package http

import (
	"net/http"

	"strconv"

	"backend_go/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ExportHandler struct {
}

func NewExportHandler(api *gin.RouterGroup) {
	handler := &ExportHandler{}

	group := api.Group("/export")
	{
		group.GET("/invoice/:id", handler.ExportInvoicePDF)
		group.GET("/medical-record/:id", handler.ExportMedicalRecordPDF)
	}
}

func (h *ExportHandler) ExportInvoicePDF(c *gin.Context) {
	// TODO: Implement actual PDF generation using gofpdf or html2pdf
	// For now, return a mock response
	patientID, _ := strconv.Atoi(c.Param("id"))
	utils.SuccessResponse(c, http.StatusOK, "Success", gin.H{
		"message": "Export PDF requested for Patient " + strconv.Itoa(int(patientID)),
		"url":     "/exports/pdf/" + strconv.Itoa(int(patientID)),
	})
}

func (h *ExportHandler) ExportMedicalRecordPDF(c *gin.Context) {
	// TODO: Implement actual PDF generation
	patientID, _ := strconv.Atoi(c.Param("id"))
	utils.SuccessResponse(c, http.StatusOK, "Success", gin.H{
		"message": "Medical Record PDF exported successfully (MOCK) for Patient " + strconv.Itoa(int(patientID)),
		"url":     "/downloads/medical_record_" + c.Param("id") + ".pdf",
	})
}
