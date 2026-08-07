package http

import (
	"net/http"
	"strconv"

	"backend_go/internal/models"
	"backend_go/internal/usecase"

	"backend_go/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TherapySessionHandler struct {
	sessionUC usecase.TherapySessionUseCase
}

func NewTherapySessionHandler(api *gin.RouterGroup, sessionUC usecase.TherapySessionUseCase) {
	handler := &TherapySessionHandler{
		sessionUC: sessionUC,
	}

	group := api.Group("/therapy-sessions")
	{
		group.GET("", handler.Fetch)
		group.GET("/:id", handler.GetByID)
		group.POST("", handler.Store)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", handler.Delete)
		group.GET("/weekly-schedule", handler.GetWeeklySchedule)
	}
}

func (h *TherapySessionHandler) Fetch(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))

	offset := (page - 1) * perPage

	sessions, total, err := h.sessionUC.Fetch(offset, perPage)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponsePaginated(c, http.StatusOK, "Success", sessions, page, perPage, total)
}

func (h *TherapySessionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	session, err := h.sessionUC.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Therapy Session not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", session)
}

func (h *TherapySessionHandler) Store(c *gin.Context) {
	var session models.TherapySession
	if err := c.ShouldBindJSON(&session); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.sessionUC.Store(&session); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success", session)
}

func (h *TherapySessionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var session models.TherapySession
	if err := c.ShouldBindJSON(&session); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.sessionUC.Update(id, &session); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Therapy Session updated successfully", nil)
}

func (h *TherapySessionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.sessionUC.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Therapy Session deleted successfully", nil)
}

func (h *TherapySessionHandler) GetWeeklySchedule(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "start_date and end_date are required", nil)
		return
	}

	sessions, err := h.sessionUC.GetWeeklySchedule(startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", sessions)
}
