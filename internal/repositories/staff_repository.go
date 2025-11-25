package repositories

import (
	"agnos-middleware/internal/models"
	"errors"

	"gorm.io/gorm"
)

type StaffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) *StaffRepository {
	return &StaffRepository{db: db}
}

func (r *StaffRepository) CreateStaff(staff *models.Staff) error {

	result := r.db.Create(staff)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *StaffRepository) GetStaffByUsername(username string) (*models.Staff, error) {
	staff := &models.Staff{}

	result := r.db.Where("username = ?", username).First(staff)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("staff not found")
		}
		return nil, result.Error
	}

	return staff, nil
}

func (r *StaffRepository) GetStaffByID(id int) (*models.Staff, error) {
	staff := &models.Staff{}

	result := r.db.First(staff, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("staff not found")
		}
		return nil, result.Error
	}

	return staff, nil
}

func (r *StaffRepository) GetStaffByEmail(email string) (*models.Staff, error) {
	staff := &models.Staff{}

	result := r.db.Where("email = ?", email).First(staff)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("staff not found")
		}
		return nil, result.Error
	}

	return staff, nil
}

func (r *StaffRepository) GetStaffByEmployeeID(employeeID string) (*models.Staff, error) {
	staff := &models.Staff{}

	result := r.db.Where("employee_id = ?", employeeID).First(staff)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("staff not found")
		}
		return nil, result.Error
	}

	return staff, nil
}
