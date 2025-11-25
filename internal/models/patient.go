package models

import (
	"time"
)

type Patient struct {
	ID           int       `json:"id" gorm:"primaryKey;column:id"`
	NationalID   *string   `json:"national_id,omitempty" gorm:"column:national_id"`
	PassportID   *string   `json:"passport_id,omitempty" gorm:"column:passport_id"`
	FirstNameTH  *string   `json:"first_name_th,omitempty" gorm:"column:first_name_th"`
	MiddleNameTH *string   `json:"middle_name_th,omitempty" gorm:"column:middle_name_th"`
	LastNameTH   *string   `json:"last_name_th,omitempty" gorm:"column:last_name_th"`
	FirstNameEN  *string   `json:"first_name_en,omitempty" gorm:"column:first_name_en"`
	MiddleNameEN *string   `json:"middle_name_en,omitempty" gorm:"column:middle_name_en"`
	LastNameEN   *string   `json:"last_name_en,omitempty" gorm:"column:last_name_en"`
	DateOfBirth  time.Time `json:"date_of_birth" gorm:"column:date_of_birth"`
	PhoneNumber  *string   `json:"phone_number,omitempty" gorm:"column:phone_number"`
	Email        *string   `json:"email,omitempty" gorm:"column:email"`
	Gender       string    `json:"gender" gorm:"column:gender"`
	PatientHN    string    `json:"patient_hn" gorm:"uniqueIndex:idx_patient_hn_hospital;column:patient_hn"`
	Hospital     string    `json:"hospital" gorm:"uniqueIndex:idx_patient_hn_hospital;column:hospital"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at"`
}

func (Patient) TableName() string {
	return "patient"
}

type PatientSearchRequest struct {
	ID          *string `form:"id" example:"1234567890123"` // Can be either national_id or passport_id (per HIS API spec)
	PatientHN   *string `form:"patient_hn"`
	NationalID  *string `form:"national_id"`
	PassportID  *string `form:"passport_id"`
	FirstName   *string `form:"first_name"`
	MiddleName  *string `form:"middle_name"`
	LastName    *string `form:"last_name"`
	DateOfBirth *string `form:"date_of_birth"`
	PhoneNumber *string `form:"phone_number"`
	Email       *string `form:"email"`
	Gender      *string `form:"gender"`
}

type PatientSearchResponse struct {
	Patients []*Patient `json:"patients,omitempty"`
	Count    int        `json:"count"`
	Error    string     `json:"error,omitempty"`
}
