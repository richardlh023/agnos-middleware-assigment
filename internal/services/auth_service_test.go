package services

import (
	"agnos-middleware/internal/configs"
	"agnos-middleware/internal/models"
	"agnos-middleware/internal/repositories"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Staff{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func TestCreateStaff_Positive(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewStaffRepository(db)
	config := &configs.ApplicationConfig{
		JWT: struct {
			Secret string
		}{Secret: "test-secret"},
	}
	service := NewAuthService(repo, config)

	req := &models.CreateStaffRequest{
		EmployeeID: "EMP001",
		Username:   "testuser",
		Password:   "password123",
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@hospital.com",
		Role:       "Doctor",
		Hospital:   "Hospital A",
	}

	staff, err := service.CreateStaff(req)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if staff.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", staff.Username)
	}

	if staff.Hospital != "Hospital A" {
		t.Errorf("Expected hospital 'Hospital A', got '%s'", staff.Hospital)
	}

	if staff.PasswordHash == "" {
		t.Error("Expected password hash to be set")
	}
}

func TestCreateStaff_Negative_DuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewStaffRepository(db)
	config := &configs.ApplicationConfig{
		JWT: struct {
			Secret string
		}{Secret: "test-secret"},
	}
	service := NewAuthService(repo, config)

	req := &models.CreateStaffRequest{
		EmployeeID: "EMP001",
		Username:   "testuser",
		Password:   "password123",
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@hospital.com",
		Role:       "Doctor",
		Hospital:   "Hospital A",
	}

	_, err := service.CreateStaff(req)
	if err != nil {
		t.Fatalf("Expected no error on first creation, got: %v", err)
	}

	_, err = service.CreateStaff(req)
	if err == nil {
		t.Error("Expected error for duplicate username, got nil")
	}

	if err.Error() != "username already exists" {
		t.Errorf("Expected 'username already exists', got: %v", err)
	}
}

func TestLogin_Positive(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewStaffRepository(db)
	config := &configs.ApplicationConfig{
		JWT: struct {
			Secret string
		}{Secret: "test-secret"},
	}
	service := NewAuthService(repo, config)

	// Create staff first
	createReq := &models.CreateStaffRequest{
		Username: "testuser",
		Password: "password123",
		Hospital: "Hospital A",
	}
	_, err := service.CreateStaff(createReq)
	if err != nil {
		t.Fatalf("Failed to create staff: %v", err)
	}

	loginReq := &models.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	response, err := service.Login(loginReq)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if response.Token == "" {
		t.Error("Expected token to be set")
	}

	if response.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", response.Username)
	}
}

func TestLogin_Negative_InvalidCredentials(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewStaffRepository(db)
	config := &configs.ApplicationConfig{
		JWT: struct {
			Secret string
		}{Secret: "test-secret"},
	}
	service := NewAuthService(repo, config)

	createReq := &models.CreateStaffRequest{
		Username: "testuser",
		Password: "password123",
		Hospital: "Hospital A",
	}
	_, err := service.CreateStaff(createReq)
	if err != nil {
		t.Fatalf("Failed to create staff: %v", err)
	}

	loginReq := &models.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	_, err = service.Login(loginReq)
	if err == nil {
		t.Error("Expected error for invalid credentials, got nil")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials', got: %v", err)
	}
}

func TestLogin_Negative_UserNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewStaffRepository(db)
	config := &configs.ApplicationConfig{
		JWT: struct {
			Secret string
		}{Secret: "test-secret"},
	}
	service := NewAuthService(repo, config)

	loginReq := &models.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	_, err := service.Login(loginReq)
	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials', got: %v", err)
	}
}

func TestValidateToken_Positive(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewStaffRepository(db)
	config := &configs.ApplicationConfig{
		JWT: struct {
			Secret string
		}{Secret: "test-secret"},
	}
	service := NewAuthService(repo, config)

	createReq := &models.CreateStaffRequest{
		Username: "testuser",
		Password: "password123",
		Hospital: "Hospital A",
	}
	_, err := service.CreateStaff(createReq)
	if err != nil {
		t.Fatalf("Failed to create staff: %v", err)
	}

	loginReq := &models.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	loginResp, err := service.Login(loginReq)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	staff, err := service.ValidateToken(loginResp.Token)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if staff.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", staff.Username)
	}
}

func TestValidateToken_Negative_InvalidToken(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewStaffRepository(db)
	config := &configs.ApplicationConfig{
		JWT: struct {
			Secret string
		}{Secret: "test-secret"},
	}
	service := NewAuthService(repo, config)

	_, err := service.ValidateToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}
