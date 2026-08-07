package http

import (
	"net/http"

	"backend_go/internal/usecase"
	"backend_go/pkg/utils"

	"github.com/gin-gonic/gin"
)

type MasterHandler struct {
	masterUC usecase.MasterUseCase
}

func NewMasterHandler(api *gin.RouterGroup, masterUC usecase.MasterUseCase) {
	handler := &MasterHandler{
		masterUC: masterUC,
	}

	api.GET("/patient-categories", handler.GetPatientCategories)
	api.GET("/genders", handler.GetGenders)
}

func (h *MasterHandler) GetPatientCategories(c *gin.Context) {
	categories, err := h.masterUC.GetPatientCategories()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", categories)
}

func (h *MasterHandler) GetGenders(c *gin.Context) {
	genders, err := h.masterUC.GetGenders()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", genders)
}
