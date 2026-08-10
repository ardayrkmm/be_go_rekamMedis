package http

import (
	"net/http"
	"strconv"

	"backend_go/internal/models"
	"backend_go/internal/usecase"
	"backend_go/internal/middleware"

	"backend_go/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

type PhysiotherapistHandler struct {
	physioUC usecase.PhysiotherapistUseCase
}

func NewPhysiotherapistHandler(api *gin.RouterGroup, physioUC usecase.PhysiotherapistUseCase) {
	handler := &PhysiotherapistHandler{
		physioUC: physioUC,
	}

	group := api.Group("/physiotherapists")
	adminOnly := middleware.RoleMiddleware(string(models.RoleAdmin), string(models.RoleOwner))
	{
		group.GET("", handler.Fetch)
		group.GET("/:id", handler.GetByID)
		group.POST("", adminOnly, handler.Store)
		group.PUT("/:id", adminOnly, handler.Update)
		group.DELETE("/:id", adminOnly, handler.Delete)
		group.POST("/:id/restore", adminOnly, handler.Restore)
	}
}

func (h *PhysiotherapistHandler) Fetch(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))
	search := c.Query("search")

	offset := (page - 1) * perPage

	physios, total, err := h.physioUC.Fetch(offset, perPage, search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponsePaginated(c, http.StatusOK, "Success", physios, page, perPage, total)
}

func (h *PhysiotherapistHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	physio, err := h.physioUC.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Physiotherapist not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", physio)
}

func (h *PhysiotherapistHandler) Store(c *gin.Context) {
	var physio models.Physiotherapist
	if err := c.ShouldBindJSON(&physio); err != nil {
		utils.HandleValidationError(c, err)
		return
	}
	
	fmt.Printf("Received Physio: %+v\n", physio)

	if err := h.physioUC.Store(&physio); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success", physio)
}

func (h *PhysiotherapistHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var physio models.Physiotherapist
	if err := c.ShouldBindJSON(&physio); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if err := h.physioUC.Update(id, &physio); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Physiotherapist updated successfully", nil)
}

func (h *PhysiotherapistHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.physioUC.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Physiotherapist deleted successfully", nil)
}

func (h *PhysiotherapistHandler) Restore(c *gin.Context) {
	id := c.Param("id")
	if err := h.physioUC.Restore(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Physiotherapist restored successfully", nil)
}
