package api

import (
	"agnos-middleware/internal/middlewares"
	"agnos-middleware/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	staffController *StaffController,
	patientController *PatientController,
	authService *services.AuthService,
) *gin.Engine {
	router := gin.Default()

	router.Use(middlewares.ErrorMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/")
	{
		api.POST("/staff/create", staffController.CreateStaff)
		api.POST("/staff/login", staffController.Login)
	}

	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware(authService))
	{
		protected.GET("/patient/search", patientController.SearchPatient)
	}

	return router
}
