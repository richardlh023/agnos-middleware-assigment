package repositories

import (
	"agnos-middleware/internal/models"
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

type PatientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) UpsertPatient(patient *models.Patient) error {
	result := r.db.Where("patient_hn = ? AND hospital = ?", patient.PatientHN, patient.Hospital).
		Assign(*patient).
		FirstOrCreate(patient)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PatientRepository) GetPatientByHN(hn string, hospital string) (*models.Patient, error) {
	patient := &models.Patient{}
	result := r.db.Where("patient_hn = ? AND hospital = ?", hn, hospital).First(patient)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, sql.ErrNoRows
		}
		return nil, result.Error
	}

	return patient, nil
}

func (r *PatientRepository) SearchPatients(req *models.PatientSearchRequest, hospital string) ([]*models.Patient, error) {
	var patients []*models.Patient
	query := r.db.Model(&models.Patient{}).Where("hospital = ?", hospital)

	if req.ID != nil && *req.ID != "" {
		// ID can be either national_id or passport_id (per HIS API spec)
		query = query.Where("(national_id = ? OR passport_id = ?)", *req.ID, *req.ID)
	}

	if req.PatientHN != nil && *req.PatientHN != "" {
		query = query.Where("patient_hn = ?", *req.PatientHN)
	}
	if req.NationalID != nil && *req.NationalID != "" {
		query = query.Where("national_id = ?", *req.NationalID)
	}
	if req.PassportID != nil && *req.PassportID != "" {
		query = query.Where("passport_id = ?", *req.PassportID)
	}

	if req.FirstName != nil && *req.FirstName != "" {
		query = query.Where("(first_name_en ILIKE ? OR first_name_th ILIKE ?)", "%"+*req.FirstName+"%", "%"+*req.FirstName+"%")
	}
	if req.MiddleName != nil && *req.MiddleName != "" {
		query = query.Where("(middle_name_en ILIKE ? OR middle_name_th ILIKE ?)", "%"+*req.MiddleName+"%", "%"+*req.MiddleName+"%")
	}
	if req.LastName != nil && *req.LastName != "" {
		query = query.Where("(last_name_en ILIKE ? OR last_name_th ILIKE ?)", "%"+*req.LastName+"%", "%"+*req.LastName+"%")
	}

	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		query = query.Where("date_of_birth = ?", *req.DateOfBirth)
	}
	if req.PhoneNumber != nil && *req.PhoneNumber != "" {
		query = query.Where("phone_number = ?", *req.PhoneNumber)
	}
	if req.Email != nil && *req.Email != "" {
		query = query.Where("email = ?", *req.Email)
	}
	if req.Gender != nil && *req.Gender != "" {
		query = query.Where("gender = ?", *req.Gender)
	}

	result := query.Find(&patients)
	if result.Error != nil {
		return nil, result.Error
	}

	return patients, nil
}
