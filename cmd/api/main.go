package main

import (
	"log"
	"os"

	_ "backend_go/docs"
	"backend_go/internal/delivery/http"
	"backend_go/internal/middleware"
	"backend_go/internal/repository"
	"backend_go/internal/usecase"
	"backend_go/pkg/logger"
	"backend_go/pkg/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Rekam Medis API
// @version 1.0
// @description API documentation for Rekam Medis
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Logger
	logger.InitLogger()
	logger.Log.Info("Starting application...")

	// Connect to database
	db := utils.ConnectDB()



	// Setup Gin router
	r := gin.Default()

	// CORS Configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Repositories
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	patientRepo := repository.NewPatientRepository(db)
	physioRepo := repository.NewPhysiotherapistRepository(db)
	serviceRepo := repository.NewServiceMasterRepository(db)
	categoryRepo := repository.NewServiceCategoryRepository(db)
	recordRepo := repository.NewMedicalRecordRepository(db)
	appointmentRepo := repository.NewAppointmentRepository(db)
	sessionRepo := repository.NewTherapySessionRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	activityLogRepo := repository.NewActivityLogRepository(db)
	masterRepo := repository.NewMasterRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)

	// Usecases
	authUC := usecase.NewAuthUseCase(userRepo, authRepo)
	userUC := usecase.NewUserUseCase(userRepo)
	patientUC := usecase.NewPatientUseCase(patientRepo)
	physioUC := usecase.NewPhysiotherapistUseCase(physioRepo, userRepo)
	serviceUC := usecase.NewServiceMasterUseCase(serviceRepo)
	categoryUC := usecase.NewServiceCategoryUseCase(categoryRepo)
	recordUC := usecase.NewMedicalRecordUseCase(recordRepo, patientRepo, physioRepo, serviceRepo, appointmentRepo)
	appointmentUC := usecase.NewAppointmentUseCase(appointmentRepo, paymentRepo, serviceRepo)
	sessionUC := usecase.NewTherapySessionUseCase(sessionRepo, paymentRepo, serviceRepo, recordRepo)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo, patientRepo, physioRepo)
	notificationUC := usecase.NewNotificationUseCase(notificationRepo)
	activityLogUC := usecase.NewActivityLogUseCase(activityLogRepo)
	masterUC := usecase.NewMasterUseCase(masterRepo)
	dashboardUC := usecase.NewDashboardUsecase(dashboardRepo)

	// Base API Group
	api := r.Group("/api/v1")

	// Public Routes
	publicAPI := api.Group("")

	// Protected Routes
	protectedAPI := api.Group("")
	protectedAPI.Use(middleware.AuthMiddleware(db))

	http.NewAuthHandler(publicAPI, protectedAPI, authUC)

	http.NewUserHandler(protectedAPI, userUC, db)
	http.NewPatientHandler(protectedAPI, patientUC)
	http.NewPhysiotherapistHandler(protectedAPI, physioUC)
	http.NewServiceMasterHandler(protectedAPI, serviceUC)
	http.NewServiceCategoryHandler(protectedAPI, categoryUC)
	http.NewAppointmentHandler(protectedAPI, appointmentUC, physioUC)
	http.NewTherapySessionHandler(protectedAPI, sessionUC, physioUC, userUC)
	http.NewMedicalRecordHandler(protectedAPI, recordUC, userUC, sessionUC, physioUC)
	http.NewPaymentHandler(protectedAPI, paymentUC)
	http.NewDashboardHandler(protectedAPI, dashboardUC)
	http.NewExportHandler(protectedAPI)
	http.NewNotificationHandler(protectedAPI, notificationUC)
	http.NewActivityLogHandler(protectedAPI, activityLogUC)
	http.NewMasterHandler(protectedAPI, masterUC)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	logger.Log.Infof("Server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		logger.Log.Fatalf("Failed to run server: %v", err)
	}
}
