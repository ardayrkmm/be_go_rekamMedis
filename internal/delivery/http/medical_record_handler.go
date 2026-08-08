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

type MedicalRecordHandler struct {
	recordUC  usecase.MedicalRecordUseCase
	userUC    usecase.UserUseCase
	sessionUC usecase.TherapySessionUseCase
	physioUC  usecase.PhysiotherapistUseCase
}

func NewMedicalRecordHandler(api *gin.RouterGroup, recordUC usecase.MedicalRecordUseCase, userUC usecase.UserUseCase, sessionUC usecase.TherapySessionUseCase, physioUC usecase.PhysiotherapistUseCase) {
	handler := &MedicalRecordHandler{
		recordUC:  recordUC,
		userUC:    userUC,
		sessionUC: sessionUC,
		physioUC:  physioUC,
	}

	group := api.Group("/medical-records")
	adminOnly := middleware.RoleMiddleware(string(models.RoleAdmin), string(models.RoleOwner))
	{
		group.GET("", handler.Fetch)
		group.GET("/:id", handler.GetByID)
		group.POST("", handler.Store)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", adminOnly, handler.Delete)
	}

	// Histori Rekam Medis Pasien
	api.GET("/patients/:id/medical-records", handler.GetHistory)
}

func (h *MedicalRecordHandler) Fetch(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))

	offset := (page - 1) * perPage

	records, total, err := h.recordUC.Fetch(offset, perPage)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	userCtx, exists := c.Get("user")
	if exists {
		user := userCtx.(models.User)
		if user.Role == string(models.RoleFisioterapis) {
			physioID := h.getPhysiotherapistID(user.Email)
			var filtered []models.MedicalRecord
			for _, rec := range records {
				if rec.PhysiotherapistID == physioID {
					filtered = append(filtered, rec)
				}
			}
			records = filtered
			total = int64(len(filtered))
		}
	}

	utils.SuccessResponsePaginated(c, http.StatusOK, "Success", records, page, perPage, total)
}

func (h *MedicalRecordHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	record, err := h.recordUC.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Medical Record not found", nil)
		return
	}

	userCtx, exists := c.Get("user")
	if exists {
		user := userCtx.(models.User)
		if user.Role == string(models.RoleFisioterapis) {
			physioID := h.getPhysiotherapistID(user.Email)
			if record.PhysiotherapistID != physioID {
				hasTreated, _ := h.sessionUC.HasTreatedPatient(physioID, record.PatientID)
				if !hasTreated {
					utils.ErrorResponse(c, http.StatusForbidden, "Akses ditolak", nil)
					return
				}
			}
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", record)
}

func (h *MedicalRecordHandler) getPhysiotherapistID(email string) string {
	physios, _, err := h.physioUC.Fetch(0, 1000)
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

func (h *MedicalRecordHandler) GetHistory(c *gin.Context) {
	patientID := c.Param("id")

	userCtx, exists := c.Get("user")
	if exists {
		user := userCtx.(models.User)
		if user.Role == string(models.RoleFisioterapis) {
			physioID := h.getPhysiotherapistID(user.Email)
			hasTreated, _ := h.sessionUC.HasTreatedPatient(physioID, patientID)
			if !hasTreated {
				utils.ErrorResponse(c, http.StatusForbidden, "Akses ditolak: Anda belum pernah menangani pasien ini", nil)
				return
			}
		}
	}

	records, err := h.recordUC.GetByPatientID(patientID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", records)
}

func (h *MedicalRecordHandler) Store(c *gin.Context) {
	var record models.MedicalRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	userCtx, exists := c.Get("user")
	if exists {
		user := userCtx.(models.User)
		if user.Role == string(models.RoleFisioterapis) {
			physioID := h.getPhysiotherapistID(user.Email)
			if record.PhysiotherapistID != physioID {
				utils.ErrorResponse(c, http.StatusForbidden, "Akses ditolak: Anda tidak bisa membuat catatan klinis untuk fisioterapis lain", nil)
				return
			}
		}
	}

	if err := h.recordUC.Store(&record); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Success", record)
}

func (h *MedicalRecordHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var record models.MedicalRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	userCtx, exists := c.Get("user")
	if exists {
		user := userCtx.(models.User)
		if user.Role == string(models.RoleFisioterapis) {
			physioID := h.getPhysiotherapistID(user.Email)
			

			existing, err := h.recordUC.GetByID(id)
			if err != nil {
				utils.ErrorResponse(c, http.StatusNotFound, "Rekam medis tidak ditemukan", nil)
				return
			}

			if existing.PhysiotherapistID != physioID || record.PhysiotherapistID != physioID {
				utils.ErrorResponse(c, http.StatusForbidden, "Akses ditolak: Anda tidak berwenang mengedit catatan klinis ini", nil)
				return
			}
		}
	}

	if err := h.recordUC.Update(id, &record); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Medical Record updated successfully", nil)
}

func (h *MedicalRecordHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	userCtx, exists := c.Get("user")
	if exists {
		user := userCtx.(models.User)
		if user.Role == string(models.RoleFisioterapis) {
			physioID := h.getPhysiotherapistID(user.Email)
			
			existing, err := h.recordUC.GetByID(id)
			if err != nil {
				utils.ErrorResponse(c, http.StatusNotFound, "Rekam medis tidak ditemukan", nil)
				return
			}

			if existing.PhysiotherapistID != physioID {
				utils.ErrorResponse(c, http.StatusForbidden, "Akses ditolak: Anda tidak berwenang menghapus catatan klinis ini", nil)
				return
			}
		}
	}

	if err := h.recordUC.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Medical Record deleted successfully", nil)
}
