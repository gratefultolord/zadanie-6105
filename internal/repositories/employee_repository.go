package repositories

import (
	"context"
	"zadanie-6105/internal/models"

	"gorm.io/gorm"
)

type EmployeeRepository interface {
	GetEmployeeIDByUsername(ctx context.Context, username string) (string, error)
	IsEmployeeExists(ctx context.Context, employeeID string) (bool, error)
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) GetEmployeeIDByUsername(ctx context.Context, username string) (string, error) {
	var employee models.Employee
	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&employee).Error
	if err != nil {
		return "", err
	}
	return employee.ID, nil
}

func (r *employeeRepository) IsEmployeeExists(ctx context.Context, employeeID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("employee").
		Where("id = ?", employeeID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
