package repositories

import (
	"context"

	"gorm.io/gorm"
)

type OrganizationRepository interface {
	IsOrganizationExists(ctx context.Context, organizationID string) (bool, error)
}

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) IsOrganizationExists(ctx context.Context, organizationID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("organization").
		Where("id = ?", organizationID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
