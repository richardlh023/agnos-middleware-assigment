package services

import (
	"agnos-middleware/internal/configs"
	"agnos-middleware/internal/models"
	"agnos-middleware/internal/repositories"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	staffRepo *repositories.StaffRepository
	config    *configs.ApplicationConfig
}

func NewAuthService(staffRepo *repositories.StaffRepository, config *configs.ApplicationConfig) *AuthService {
	return &AuthService{
		staffRepo: staffRepo,
		config:    config,
	}
}

func (s *AuthService) CreateStaff(req *models.CreateStaffRequest) (*models.Staff, error) {
	existing, _ := s.staffRepo.GetStaffByUsername(req.Username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	existingByEmail, _ := s.staffRepo.GetStaffByEmail(req.Email)
	if existingByEmail != nil {
		return nil, errors.New("email already exists")
	}

	existingByEmployeeID, _ := s.staffRepo.GetStaffByEmployeeID(req.EmployeeID)
	if existingByEmployeeID != nil {
		return nil, errors.New("employee_id already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	staff := &models.Staff{
		EmployeeID:   req.EmployeeID,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		Role:         req.Role,
		Department:   req.Department,
		Hospital:     req.Hospital,
		IsActive:     true,
	}

	if err := s.staffRepo.CreateStaff(staff); err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	staff, err := s.staffRepo.GetStaffByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateJWT(staff)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	department := ""
	if staff.Department != nil {
		department = *staff.Department
	}

	return &models.LoginResponse{
		Token:      token,
		EmployeeID: staff.EmployeeID,
		Username:   staff.Username,
		FirstName:  staff.FirstName,
		LastName:   staff.LastName,
		Email:      staff.Email,
		Role:       staff.Role,
		Department: department,
		Hospital:   staff.Hospital,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*models.Staff, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		staffID := int(claims["staff_id"].(float64))

		staff, err := s.staffRepo.GetStaffByID(staffID)
		if err != nil {
			return nil, err
		}
		return staff, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) generateJWT(staff *models.Staff) (string, error) {
	claims := jwt.MapClaims{
		"staff_id": staff.ID,
		"username": staff.Username,
		"hospital": staff.Hospital,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.config.JWT.Secret))
}
