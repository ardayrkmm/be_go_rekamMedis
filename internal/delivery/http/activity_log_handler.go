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

type ActivityLogHandler struct {
	logUC usecase.ActivityLogUseCase
}

func NewActivityLogHandler(api *gin.RouterGroup, logUC usecase.ActivityLogUseCase) {
	handler := &ActivityLogHandler{
		logUC: logUC,
	}

	group := api.Group("/activity-logs")
	adminOnly := middleware.RoleMiddleware(string(models.RoleAdmin), string(models.RoleOwner))
	group.Use(adminOnly)
	{
		group.GET("", handler.Fetch)
		group.GET("/:id", handler.GetByID)
		group.DELETE("/:id", handler.Delete)
	}
}

func (h *ActivityLogHandler) Fetch(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))

	offset := (page - 1) * perPage

	logs, total, err := h.logUC.Fetch(offset, perPage)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponsePaginated(c, http.StatusOK, "Success", logs, page, perPage, total)
}

func (h *ActivityLogHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	log, err := h.logUC.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Activity Log not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", log)
}

func (h *ActivityLogHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.logUC.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Activity Log deleted successfully", nil)
}
