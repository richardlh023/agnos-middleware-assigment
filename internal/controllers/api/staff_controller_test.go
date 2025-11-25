package api

import (
	"agnos-middleware/internal/configs"
	"agnos-middleware/internal/models"
	"agnos-middleware/internal/repositories"
	"agnos-middleware/internal/services"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *services.AuthService) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&models.Staff{})

	repo := repositories.NewStaffRepository(db)
	config := &configs.ApplicationConfig{
		JWT: struct {
			Secret string
		}{Secret: "test-secret"},
	}
	authService := services.NewAuthService(repo, config)
	staffController := NewStaffController(authService)

	router := gin.New()
	router.POST("/staff/create", staffController.CreateStaff)
	router.POST("/staff/login", staffController.Login)

	return router, authService
}

func TestCreateStaff_Positive(t *testing.T) {
	router, _ := setupTestRouter(t)

	payload := models.CreateStaffRequest{
		EmployeeID: "EMP001",
		Username:   "testuser",
		Password:   "password123",
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@hospital.com",
		Role:       "Doctor",
		Hospital:   "Hospital A",
	}
	jsonValue, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "testuser", response["username"])
	assert.Equal(t, "Hospital A", response["hospital"])
}

func TestCreateStaff_Negative_DuplicateUsername(t *testing.T) {
	router, _ := setupTestRouter(t)

	payload := models.CreateStaffRequest{
		EmployeeID: "EMP001",
		Username:   "testuser",
		Password:   "password123",
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@hospital.com",
		Role:       "Doctor",
		Hospital:   "Hospital A",
	}
	jsonValue, _ := json.Marshal(payload)

	req1, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonValue))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	req2, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonValue))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestLogin_Positive(t *testing.T) {
	router, _ := setupTestRouter(t)

	createPayload := models.CreateStaffRequest{
		EmployeeID: "EMP001",
		Username:   "testuser",
		Password:   "password123",
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@hospital.com",
		Role:       "Doctor",
		Hospital:   "Hospital A",
	}
	createJson, _ := json.Marshal(createPayload)
	createReq, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(createJson))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	loginPayload := models.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	loginJson, _ := json.Marshal(loginPayload)
	loginReq, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(loginJson))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)

	assert.Equal(t, http.StatusOK, loginW.Code)

	var response models.LoginResponse
	json.Unmarshal(loginW.Body.Bytes(), &response)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "EMP001", response.EmployeeID)
	assert.Equal(t, "testuser", response.Username)
	assert.Equal(t, "John", response.FirstName)
	assert.Equal(t, "Doe", response.LastName)
	assert.Equal(t, "john.doe@hospital.com", response.Email)
	assert.Equal(t, "Doctor", response.Role)
	assert.Equal(t, "Hospital A", response.Hospital)
}

func TestLogin_Negative_InvalidCredentials(t *testing.T) {
	router, _ := setupTestRouter(t)

	createPayload := models.CreateStaffRequest{
		EmployeeID: "EMP001",
		Username:   "testuser",
		Password:   "password123",
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@hospital.com",
		Role:       "Doctor",
		Hospital:   "Hospital A",
	}
	createJson, _ := json.Marshal(createPayload)
	createReq, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(createJson))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	loginPayload := models.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	loginJson, _ := json.Marshal(loginPayload)
	loginReq, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(loginJson))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)

	assert.Equal(t, http.StatusUnauthorized, loginW.Code)
}
