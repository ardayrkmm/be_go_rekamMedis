package http

import (
	"net/http"
	"strconv"
	"strings"

	"backend_go/internal/models"
	"backend_go/internal/usecase"
	"backend_go/internal/middleware"

	"backend_go/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TherapySessionHandler struct {
	sessionUC usecase.TherapySessionUseCase
	physioUC  usecase.PhysiotherapistUseCase
	userUC    usecase.UserUseCase
}

func NewTherapySessionHandler(api *gin.RouterGroup, sessionUC usecase.TherapySessionUseCase, physioUC usecase.PhysiotherapistUseCase, userUC usecase.UserUseCase) {
	handler := &TherapySessionHandler{
		sessionUC: sessionUC,
		physioUC:  physioUC,
		userUC:    userUC,
	}

	group := api.Group("/therapy-sessions")
	adminOnly := middleware.RoleMiddleware(string(models.RoleAdmin), string(models.RoleOwner))
	{
		group.GET("", handler.Fetch)
		group.GET("/:id", handler.GetByID)
		group.POST("", adminOnly, handler.Store)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", adminOnly, handler.Delete)
		group.GET("/weekly-schedule", handler.GetWeeklySchedule)
	}
}

func (h *TherapySessionHandler) Fetch(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))
	patientID := c.Query("patient_id")

	offset := (page - 1) * perPage

	// GLOBAL VISIBILITY: semua role melihat seluruh sesi klinik.
	// Authorization edit dilakukan saat update/delete.
	sessions, total, err := h.sessionUC.Fetch(offset, perPage, patientID)
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

	// GLOBAL VISIBILITY: semua role yang terautentikasi bisa membaca detail sesi.
	// Fisioterapis bisa melihat sesi klinik (termasuk sesi orang lain) agar kalender akurat.
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

	existing, err := h.sessionUC.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Session not found", nil)
		return
	}

	role, _ := c.Get("role")
	if role == string(models.RoleFisioterapis) {
		userID, _ := c.Get("userID")
		user, err := h.userUC.GetByID(userID.(string))
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}
		
		physioID := h.getPhysiotherapistID(user.Email)
		if existing.PhysiotherapistID != physioID {
			utils.ErrorResponse(c, http.StatusForbidden, "Akses ditolak: Anda tidak bisa mengubah sesi orang lain", nil)
			return
		}
	}

	var session models.TherapySession
	if err := c.ShouldBindJSON(&session); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Validate status transition
	if session.Status != "" && session.Status != existing.Status {
		valid := false
		oldStatus := strings.ToLower(existing.Status)
		newStatus := strings.ToLower(session.Status)

		switch oldStatus {
		case "scheduled":
			if newStatus == "telah_tiba" || newStatus == "patient arrived" || newStatus == "cancelled" || newStatus == "rescheduled" {
				valid = true
			}
		case "telah_tiba", "patient arrived":
			if newStatus == "ongoing" || newStatus == "in progress" || newStatus == "cancelled" {
				valid = true
			}
		case "ongoing", "in progress":
			if newStatus == "completed" {
				valid = true
			}
		}
		
		if !valid {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid status transition from "+existing.Status+" to "+session.Status, nil)
			return
		}
	}

	if err := h.sessionUC.Update(id, &session); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	updated, _ := h.sessionUC.GetByID(id)
	utils.SuccessResponse(c, http.StatusOK, "Therapy Session updated successfully", updated)
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

func (h *TherapySessionHandler) getPhysiotherapistID(email string) string {
	physios, _, err := h.physioUC.Fetch(0, 1000, "")
	if err != nil {
		return ""
	}
	for _, p := range physios {
		if p.Email == email {
			return p.ID
		}
	}
	return ""
}
