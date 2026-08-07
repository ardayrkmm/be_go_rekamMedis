package http

import (
	"net/http"
	"strconv"

	"backend_go/internal/models"
	"backend_go/internal/usecase"

	"backend_go/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	appointmentUC usecase.AppointmentUseCase
}

func NewAppointmentHandler(api *gin.RouterGroup, appointmentUC usecase.AppointmentUseCase) {
	handler := &AppointmentHandler{
		appointmentUC: appointmentUC,
	}

	group := api.Group("/appointments")
	{
		group.GET("", handler.Fetch)
		group.GET("/:id", handler.GetByID)
		group.POST("", handler.Store)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", handler.Delete)
		group.POST("/:id/cancel", handler.Cancel)
		group.POST("/:id/reschedule", handler.Reschedule)
	}
}

func (h *AppointmentHandler) Fetch(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))

	offset := (page - 1) * perPage

	appointments, total, err := h.appointmentUC.Fetch(offset, perPage)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponsePaginated(c, http.StatusOK, "Success", appointments, page, perPage, total)
}

func (h *AppointmentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	appointment, err := h.appointmentUC.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Appointment not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", appointment)
}

func (h *AppointmentHandler) Store(c *gin.Context) {
	var appointment models.Appointment
	if err := c.ShouldBindJSON(&appointment); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.appointmentUC.Store(&appointment); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success", appointment)
}

func (h *AppointmentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var appointment models.Appointment
	if err := c.ShouldBindJSON(&appointment); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.appointmentUC.Update(id, &appointment); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Appointment updated successfully", nil)
}

func (h *AppointmentHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.appointmentUC.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Appointment deleted successfully", nil)
}

func (h *AppointmentHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	if err := h.appointmentUC.Cancel(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Appointment cancelled successfully", nil)
}

func (h *AppointmentHandler) Reschedule(c *gin.Context) {
	id := c.Param("id")
	var req models.Appointment
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.appointmentUC.Reschedule(id, &req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Appointment rescheduled successfully", nil)
}
