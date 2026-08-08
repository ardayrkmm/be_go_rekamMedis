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

type PaymentHandler struct {
	paymentUC usecase.PaymentUseCase
}

func NewPaymentHandler(api *gin.RouterGroup, paymentUC usecase.PaymentUseCase) {
	handler := &PaymentHandler{
		paymentUC: paymentUC,
	}

	group := api.Group("/payments")
	group.Use(middleware.RoleMiddleware(string(models.RoleAdmin), string(models.RoleOwner)))
	{
		group.GET("", handler.Fetch)
		group.GET("/:id", handler.GetByID)
		group.POST("", handler.Store)
		group.PUT("/:id", handler.Update)
		group.PATCH("/:id/status", handler.UpdateStatus)
		group.GET("/:id/pdf/download", handler.DownloadPDF)
		group.GET("/:id/pdf/receipt/download", handler.DownloadReceiptPDF)
	}
}

func (h *PaymentHandler) Fetch(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))
	search := c.Query("search")
	status := c.Query("status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	offset := (page - 1) * perPage

	payments, total, err := h.paymentUC.Fetch(offset, perPage, search, status, startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponsePaginated(c, http.StatusOK, "Success", payments, page, perPage, total)
}

func (h *PaymentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	payment, err := h.paymentUC.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Payment not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", payment)
}

func (h *PaymentHandler) Store(c *gin.Context) {
	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.paymentUC.Store(&payment); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success", payment)
}

func (h *PaymentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req models.Payment
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.paymentUC.Update(id, &req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment updated successfully", req)
}

type updateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *PaymentHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.paymentUC.UpdateStatus(id, req.Status); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment status updated successfully", nil)
}

func (h *PaymentHandler) DownloadPDF(c *gin.Context) {
	// Mock PDF bytes
	pdfBytes := []byte("%PDF-1.4\n%EOF\n")
	
	c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (h *PaymentHandler) DownloadReceiptPDF(c *gin.Context) {
	// Mock PDF bytes
	pdfBytes := []byte("%PDF-1.4\n%EOF\n")
	
	c.Header("Content-Disposition", "attachment; filename=receipt.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
