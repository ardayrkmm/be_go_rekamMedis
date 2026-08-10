package http

import (
	"net/http"
	"backend_go/internal/usecase"
	"backend_go/internal/models"
	"backend_go/internal/middleware"
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
	adminOnly := middleware.RoleMiddleware(string(models.RoleAdmin), string(models.RoleOwner))
	fisioOnly := middleware.RoleMiddleware(string(models.RoleFisioterapis))
	{
		group.GET("/admin", adminOnly, handler.GetAdminDashboard)
		group.GET("/fisio", fisioOnly, handler.GetFisioDashboard)
	}

	reportGroup := api.Group("/reports")
	{
		reportGroup.GET("/dashboard", adminOnly, handler.GetReportsDashboard)
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
	userID := c.GetString("userID")
	data, err := h.usecase.GetFisioDashboardData(userID)
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
