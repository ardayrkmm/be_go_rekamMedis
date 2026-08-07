package http

import (
	"net/http"
	"strconv"

	"backend_go/internal/usecase"
	"backend_go/pkg/utils"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationUC usecase.NotificationUseCase
}

func NewNotificationHandler(api *gin.RouterGroup, notificationUC usecase.NotificationUseCase) {
	handler := &NotificationHandler{
		notificationUC: notificationUC,
	}

	group := api.Group("/notifications")
	{
		group.GET("", handler.Fetch)
		group.GET("/unread", handler.FetchUnread)
		group.POST("/read-all", handler.MarkAllAsRead)
		group.POST("/:id/read", handler.MarkAsRead)
		group.POST("/:id/unread", handler.MarkAsUnread)
		group.DELETE("/:id", handler.Delete)
	}
}

func (h *NotificationHandler) Fetch(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))

	offset := (page - 1) * perPage

	notifications, total, err := h.notificationUC.Fetch(offset, perPage)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponsePaginated(c, http.StatusOK, "Success", notifications, page, perPage, total)
}

func (h *NotificationHandler) FetchUnread(c *gin.Context) {
	notifications, err := h.notificationUC.FetchUnread()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", notifications)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	if err := h.notificationUC.MarkAllAsRead(); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "All notifications marked as read", nil)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	if err := h.notificationUC.MarkAsRead(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Notification marked as read", nil)
}

func (h *NotificationHandler) MarkAsUnread(c *gin.Context) {
	id := c.Param("id")
	if err := h.notificationUC.MarkAsUnread(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Notification marked as unread", nil)
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.notificationUC.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Notification deleted", nil)
}
