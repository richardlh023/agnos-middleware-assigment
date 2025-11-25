package services

import (
	"agnos-middleware/internal/configs"
	"agnos-middleware/internal/models"
	"agnos-middleware/internal/repositories"
	"agnos-middleware/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PatientService struct {
	patientRepo *repositories.PatientRepository
	config      *configs.ApplicationConfig
	httpClient  *http.Client
}

func NewPatientService(patientRepo *repositories.PatientRepository, config *configs.ApplicationConfig) *PatientService {
	return &PatientService{
		patientRepo: patientRepo,
		config:      config,
		httpClient:  utils.DefaultHTTPClient(),
	}
}

func (s *PatientService) SearchPatient(req *models.PatientSearchRequest, staffHospital string) ([]*models.Patient, error) {
	patients, err := s.patientRepo.SearchPatients(req, staffHospital)
	if err != nil {
		return nil, err
	}

	if len(patients) > 0 {
		return patients, nil
	}

	if req.ID != nil && *req.ID != "" {
		patient, err := s.searchPatientFromHIS(*req.ID)
		if err != nil {
			return []*models.Patient{}, nil
		}

		if patient.Hospital != staffHospital {
			return nil, errors.New("access denied: patient does not belong to your hospital")
		}

		if err := s.patientRepo.UpsertPatient(patient); err != nil {
		}

		return []*models.Patient{patient}, nil
	}

	return []*models.Patient{}, nil
}

func (s *PatientService) searchPatientFromHIS(patientID string) (*models.Patient, error) {
	fmt.Printf("[HIS API] Searching for patient: %s\n", patientID)

	// Try to call actual API first
	patient, err := s.callHISAPI(patientID)
	if err == nil && patient != nil {
		fmt.Printf("[HIS API] Successfully fetched from external API\n")
		return patient, nil
	}

	// If API call failed, fall back to mock data since current cant connect to api url
	fmt.Printf("[HIS API] API call failed, using mock data: %v\n", err)
	patient = s.getMockPatient(patientID)
	if patient != nil {
		fmt.Printf("[HIS API] Found in mock data\n")
		return patient, nil
	}

	return nil, errors.New("patient not found in HIS API")
}

func (s *PatientService) callHISAPI(patientID string) (*models.Patient, error) {
	url := fmt.Sprintf("%s/patient/search/%s", s.config.HISAPI.BaseURL, patientID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HIS API returned status %d: %s", resp.StatusCode, string(body))
	}

	var patient models.Patient
	if err := json.NewDecoder(resp.Body).Decode(&patient); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &patient, nil
}

func (s *PatientService) getMockPatient(patientID string) *models.Patient {
	mockPatients := map[string]*models.Patient{
		"HN001": {
			NationalID:   stringPtr("1234567890123"),
			PassportID:   nil,
			FirstNameTH:  stringPtr("สมชาย"),
			MiddleNameTH: nil,
			LastNameTH:   stringPtr("ใจดี"),
			FirstNameEN:  stringPtr("Somchai"),
			MiddleNameEN: nil,
			LastNameEN:   stringPtr("Jaidee"),
			DateOfBirth:  time.Date(1985, 3, 15, 0, 0, 0, 0, time.UTC),
			PhoneNumber:  stringPtr("0891234567"),
			Email:        stringPtr("somchai@email.com"),
			Gender:       "M",
			PatientHN:    "HN001",
			Hospital:     "Hospital A",
		},
		"HN002": {
			NationalID:   stringPtr("9876543210987"),
			PassportID:   nil,
			FirstNameTH:  stringPtr("สมหญิง"),
			MiddleNameTH: nil,
			LastNameTH:   stringPtr("รักดี"),
			FirstNameEN:  stringPtr("Somying"),
			MiddleNameEN: nil,
			LastNameEN:   stringPtr("Rakdee"),
			DateOfBirth:  time.Date(1990, 7, 20, 0, 0, 0, 0, time.UTC),
			PhoneNumber:  stringPtr("0899876543"),
			Email:        stringPtr("somying@email.com"),
			Gender:       "F",
			PatientHN:    "HN002",
			Hospital:     "Hospital A",
		},
		"HN003": {
			NationalID:   nil,
			PassportID:   stringPtr("AB1234567"),
			FirstNameTH:  nil,
			MiddleNameTH: nil,
			LastNameTH:   nil,
			FirstNameEN:  stringPtr("John"),
			MiddleNameEN: stringPtr("William"),
			LastNameEN:   stringPtr("Smith"),
			DateOfBirth:  time.Date(1978, 11, 5, 0, 0, 0, 0, time.UTC),
			PhoneNumber:  stringPtr("+1234567890"),
			Email:        stringPtr("john.smith@email.com"),
			Gender:       "M",
			PatientHN:    "HN003",
			Hospital:     "Hospital A",
		},
		"HN005": {
			NationalID:   stringPtr("1122334455667"),
			PassportID:   nil,
			FirstNameTH:  stringPtr("ประเสริฐ"),
			MiddleNameTH: stringPtr("สุข"),
			LastNameTH:   stringPtr("สมบูรณ์"),
			FirstNameEN:  stringPtr("Prasert"),
			MiddleNameEN: stringPtr("Suk"),
			LastNameEN:   stringPtr("Sombun"),
			DateOfBirth:  time.Date(1992, 2, 14, 0, 0, 0, 0, time.UTC),
			PhoneNumber:  stringPtr("0823456789"),
			Email:        stringPtr("prasert@email.com"),
			Gender:       "M",
			PatientHN:    "HN005",
			Hospital:     "Hospital A",
		},
		"HN006": {
			NationalID:   stringPtr("2233445566778"),
			PassportID:   nil,
			FirstNameTH:  stringPtr("มาลี"),
			MiddleNameTH: nil,
			LastNameTH:   stringPtr("ดีใจ"),
			FirstNameEN:  stringPtr("Malee"),
			MiddleNameEN: nil,
			LastNameEN:   stringPtr("Deejai"),
			DateOfBirth:  time.Date(1988, 9, 30, 0, 0, 0, 0, time.UTC),
			PhoneNumber:  stringPtr("0834567890"),
			Email:        stringPtr("malee@email.com"),
			Gender:       "F",
			PatientHN:    "HN006",
			Hospital:     "Hospital A",
		},
		"HN007": {
			NationalID:   stringPtr("3344556677889"),
			PassportID:   nil,
			FirstNameTH:  stringPtr("สมศักดิ์"),
			MiddleNameTH: nil,
			LastNameTH:   stringPtr("เก่งดี"),
			FirstNameEN:  stringPtr("Somsak"),
			MiddleNameEN: nil,
			LastNameEN:   stringPtr("Kengdee"),
			DateOfBirth:  time.Date(1980, 6, 25, 0, 0, 0, 0, time.UTC),
			PhoneNumber:  stringPtr("0845678901"),
			Email:        nil,
			Gender:       "M",
			PatientHN:    "HN007",
			Hospital:     "Hospital A",
		},
		"HN008": {
			NationalID:   nil,
			PassportID:   stringPtr("CD9876543"),
			FirstNameTH:  nil,
			MiddleNameTH: nil,
			LastNameTH:   nil,
			FirstNameEN:  stringPtr("Sarah"),
			MiddleNameEN: stringPtr("Jane"),
			LastNameEN:   stringPtr("Johnson"),
			DateOfBirth:  time.Date(1995, 4, 12, 0, 0, 0, 0, time.UTC),
			PhoneNumber:  stringPtr("+44123456789"),
			Email:        stringPtr("sarah.j@email.com"),
			Gender:       "F",
			PatientHN:    "HN008",
			Hospital:     "Hospital A",
		},
		"HN004": {
			NationalID:   stringPtr("1111222233334"),
			PassportID:   nil,
			FirstNameTH:  stringPtr("วิชัย"),
			MiddleNameTH: nil,
			LastNameTH:   stringPtr("สุขใจ"),
			FirstNameEN:  stringPtr("Wichai"),
			MiddleNameEN: nil,
			LastNameEN:   stringPtr("Sukjai"),
			DateOfBirth:  time.Date(1982, 5, 10, 0, 0, 0, 0, time.UTC),
			PhoneNumber:  stringPtr("0812345678"),
			Email:        nil,
			Gender:       "M",
			PatientHN:    "HN004",
			Hospital:     "Hospital B",
		},
		"HN009": {
			NationalID:   stringPtr("4455667788990"),
			PassportID:   nil,
			FirstNameTH:  stringPtr("นิดา"),
			MiddleNameTH: nil,
			LastNameTH:   stringPtr("รุ่งเรือง"),
			FirstNameEN:  stringPtr("Nida"),
			MiddleNameEN: nil,
			LastNameEN:   stringPtr("Rungruang"),
			DateOfBirth:  time.Date(1993, 8, 18, 0, 0, 0, 0, time.UTC),
			PhoneNumber:  stringPtr("0856789012"),
			Email:        stringPtr("nida@email.com"),
			Gender:       "F",
			PatientHN:    "HN009",
			Hospital:     "Hospital B",
		},
		"HN010": {
			NationalID:   stringPtr("5566778899001"),
			PassportID:   nil,
			FirstNameTH:  stringPtr("วีระ"),
			MiddleNameTH: stringPtr("ชัย"),
			LastNameTH:   stringPtr("วัฒนา"),
			FirstNameEN:  stringPtr("Weera"),
			MiddleNameEN: stringPtr("Chai"),
			LastNameEN:   stringPtr("Wattana"),
			DateOfBirth:  time.Date(1987, 12, 3, 0, 0, 0, 0, time.UTC),
			PhoneNumber:  stringPtr("0867890123"),
			Email:        stringPtr("weera@email.com"),
			Gender:       "M",
			PatientHN:    "HN010",
			Hospital:     "Hospital B",
		},
		"HN011": {
			NationalID:   stringPtr("6677889900112"),
			PassportID:   nil,
			FirstNameTH:  stringPtr("สุภาพ"),
			MiddleNameTH: nil,
			LastNameTH:   stringPtr("ใจดี"),
			FirstNameEN:  stringPtr("Supap"),
			MiddleNameEN: nil,
			LastNameEN:   stringPtr("Jaidee"),
			DateOfBirth:  time.Date(1991, 1, 22, 0, 0, 0, 0, time.UTC),
			PhoneNumber:  stringPtr("0878901234"),
			Email:        nil,
			Gender:       "F",
			PatientHN:    "HN011",
			Hospital:     "Hospital B",
		},
	}

	for _, p := range mockPatients {
		if p.NationalID != nil && *p.NationalID == patientID {
			return p
		}
		if p.PassportID != nil && *p.PassportID == patientID {
			return p
		}
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}
