package services

import (
	"agnos-middleware/internal/configs"
	"agnos-middleware/internal/models"
	"agnos-middleware/internal/repositories"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupPatientTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Patient{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func getTestConfig() *configs.ApplicationConfig {
	return &configs.ApplicationConfig{
		HISAPI: struct {
			BaseURL string
		}{
			BaseURL: "https://hospital-a.api.co.th",
		},
	}
}

func TestSearchPatient_Positive_FoundInDB(t *testing.T) {
	db := setupPatientTestDB(t)
	repo := repositories.NewPatientRepository(db)
	config := getTestConfig()
	service := NewPatientService(repo, config)

	patient := &models.Patient{
		PatientHN:   "HN001",
		Hospital:    "Hospital A",
		FirstNameEN: stringPtr("John"),
		LastNameEN:  stringPtr("Doe"),
		Gender:      "M",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	err := repo.UpsertPatient(patient)
	if err != nil {
		t.Fatalf("Failed to create patient: %v", err)
	}

	req := &models.PatientSearchRequest{
		PatientHN: stringPtr("HN001"), // Use patient_hn for this test since patient is in DB
	}

	patients, err := service.SearchPatient(req, "Hospital A")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(patients) != 1 {
		t.Fatalf("Expected 1 patient, got %d", len(patients))
	}

	if patients[0].PatientHN != "HN001" {
		t.Errorf("Expected PatientHN 'HN001', got '%s'", patients[0].PatientHN)
	}
}

func TestSearchPatient_Positive_FoundInHIS(t *testing.T) {
	db := setupPatientTestDB(t)
	repo := repositories.NewPatientRepository(db)
	config := getTestConfig()
	service := NewPatientService(repo, config)

	req := &models.PatientSearchRequest{
		ID: stringPtr("9876543210987"), // Use national_id (matches HN002's national_id)
	}

	patients, err := service.SearchPatient(req, "Hospital A")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(patients) != 1 {
		t.Fatalf("Expected 1 patient, got %d", len(patients))
	}

	dbPatient, err := repo.GetPatientByHN("HN002", "Hospital A")
	if err != nil {
		t.Fatalf("Expected patient to be cached, got error: %v", err)
	}

	if dbPatient.PatientHN != "HN002" {
		t.Errorf("Expected cached PatientHN 'HN002', got '%s'", dbPatient.PatientHN)
	}
}

func TestSearchPatient_Negative_AccessDenied(t *testing.T) {
	db := setupPatientTestDB(t)
	repo := repositories.NewPatientRepository(db)
	config := getTestConfig()
	service := NewPatientService(repo, config)

	req := &models.PatientSearchRequest{
		ID: stringPtr("1111222233334"), // Use national_id instead of patient_hn
	}

	_, err := service.SearchPatient(req, "Hospital A")
	if err == nil {
		t.Error("Expected error for access denied, got nil")
	}

	if err.Error() != "access denied: patient does not belong to your hospital" {
		t.Errorf("Expected 'access denied', got: %v", err)
	}
}

func TestSearchPatient_Negative_NotFound(t *testing.T) {
	db := setupPatientTestDB(t)
	repo := repositories.NewPatientRepository(db)
	config := getTestConfig()
	service := NewPatientService(repo, config)

	req := &models.PatientSearchRequest{
		ID: stringPtr("9999999999999"), // Use a non-existent national_id
	}

	patients, err := service.SearchPatient(req, "Hospital A")
	if err != nil {
		t.Fatalf("Expected no error (empty result), got: %v", err)
	}

	if len(patients) != 0 {
		t.Fatalf("Expected 0 patients, got %d", len(patients))
	}
}
