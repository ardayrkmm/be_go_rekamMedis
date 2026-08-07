package http

import (
	"net/http"
	"backend_go/internal/usecase"
	"backend_go/pkg/utils"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	usecase usecase.DashboardUsecase
}

func NewDashboardHandler(api *gin.RouterGroup, uc usecase.DashboardUsecase) {
	handler := &DashboardHandler{
		usecase: uc,
	}

	group := api.Group("/dashboard")
	{
		group.GET("/admin", handler.GetAdminDashboard)
		group.GET("/fisio", handler.GetFisioDashboard)
	}

	reportGroup := api.Group("/reports")
	{
		reportGroup.GET("/dashboard", handler.GetReportsDashboard)
	}
}

func (h *DashboardHandler) GetAdminDashboard(c *gin.Context) {
	filter := c.Query("filter")
	data, err := h.usecase.GetAdminDashboardData(filter)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Success", data)
}

func (h *DashboardHandler) GetFisioDashboard(c *gin.Context) {
	// Ideally we get physiotherapist ID from token, but for now just pass empty or a dummy ID
	physiotherapistID := "" 
	data, err := h.usecase.GetFisioDashboardData(physiotherapistID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Success", data)
}

func (h *DashboardHandler) GetReportsDashboard(c *gin.Context) {
	data, err := h.usecase.GetReportsDashboardData()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Success", data)
}
