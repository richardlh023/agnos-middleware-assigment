package models

import (
	"time"
)

type Staff struct {
	ID           int       `json:"id" gorm:"primaryKey;column:id"`
	EmployeeID   string    `json:"employee_id" gorm:"uniqueIndex;column:employee_id"`
	Username     string    `json:"username" gorm:"uniqueIndex;column:username"`
	PasswordHash string    `json:"-" gorm:"column:password_hash"`
	FirstName    string    `json:"first_name" gorm:"column:first_name"`
	LastName     string    `json:"last_name" gorm:"column:last_name"`
	Email        string    `json:"email" gorm:"uniqueIndex;column:email"`
	PhoneNumber  *string   `json:"phone_number,omitempty" gorm:"column:phone_number"`
	Role         string    `json:"role" gorm:"column:role"`
	Department   *string   `json:"department,omitempty" gorm:"column:department"`
	Hospital     string    `json:"hospital" gorm:"column:hospital"`
	IsActive     bool      `json:"is_active" gorm:"default:true;column:is_active"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime;column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at"`
}

func (Staff) TableName() string {
	return "staff"
}

type CreateStaffRequest struct {
	EmployeeID  string  `json:"employee_id" binding:"required" example:"EMP001"`
	Username    string  `json:"username" binding:"required" example:"doctor1"`
	Password    string  `json:"password" binding:"required,min=6" example:"password123"`
	FirstName   string  `json:"first_name" binding:"required" example:"John"`
	LastName    string  `json:"last_name" binding:"required" example:"Doe"`
	Email       string  `json:"email" binding:"required,email" example:"john.doe@hospital.com"`
	PhoneNumber *string `json:"phone_number,omitempty" example:"0891234567"`
	Role        string  `json:"role" binding:"required" example:"Doctor"`
	Department  *string `json:"department,omitempty" example:"Cardiology"`
	Hospital    string  `json:"hospital" binding:"required" example:"Hospital A"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"doctor1"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type LoginResponse struct {
	Token      string `json:"token"`
	EmployeeID string `json:"employee_id"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department,omitempty"`
	Hospital   string `json:"hospital"`
}
