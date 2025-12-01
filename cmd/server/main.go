package main

// @title           Agnos Middleware API
// @version         1.0
// @description     Hospital Middleware System API for staff authentication and patient search

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

import (
	"agnos-middleware/internal/configs"
	"agnos-middleware/internal/controllers/api"
	"agnos-middleware/internal/repositories"
	"agnos-middleware/internal/services"
	"agnos-middleware/internal/utils"
	"fmt"
	"log"

	_ "agnos-middleware/docs"
)

func main() {
	fmt.Println("=====================================")

	config := configs.LoadConfig()
	fmt.Println("Configuration loaded")

	db, err := utils.ConnectDatabase(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Database connected")

	staffRepo := repositories.NewStaffRepository(db)
	patientRepo := repositories.NewPatientRepository(db)
	fmt.Println("Repositories initialized")

	authService := services.NewAuthService(staffRepo, config)
	patientService := services.NewPatientService(patientRepo, config)
	fmt.Println("Services initialized")

	staffController := api.NewStaffController(authService)
	patientController := api.NewPatientController(patientService)
	fmt.Println("Controllers initialized")

	router := api.SetupRouter(staffController, patientController, authService)
	fmt.Println("Routes configured")

	port := config.App.Port
	fmt.Printf("\n Server running on port %s\n", port)
	fmt.Printf(" Health check: http://localhost:%s/health\n", port)
	fmt.Printf(" Swagger UI: http://localhost:%s/swagger/index.html\n", port)
	fmt.Printf(" Create staff: POST http://localhost:%s/staff/create\n", port)
	fmt.Printf(" Login: POST http://localhost:%s/staff/login\n", port)
	fmt.Printf(" Search patient: GET http://localhost:%s/patient/search?id=HN001\n", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
