package repositories

import (
	"context"
	"database/sql"
	"zadanie-6105/internal/models"

	"gorm.io/gorm"
)

type TenderRepository interface {
	CreateTender(ctx context.Context, tender *models.Tender) error
	GetTenderByID(ctx context.Context, id string) (*models.Tender, error)
	GetTendersByUser(ctx context.Context, username string, limit, offset int) ([]*models.Tender, error)
	GetTenders(ctx context.Context, serviceTypes []models.TenderServiceType, limit, offset int) ([]*models.Tender, error)
	UpdateTenderStatus(ctx context.Context, id string, status models.TenderStatus) error
	UpdateTender(ctx context.Context, tender *models.Tender) error
	DeleteTender(ctx context.Context, id string) error
	GetTenderVersions(ctx context.Context, id string) ([]*models.Tender, error)
	// RollbackTenderVersion(ctx context.Context, id string, version int) error
	IsUserResponsibleForOrganization(username, organizationID string) (bool, error)
	IsUserResponsibleForTender(username, tenderId string) (bool, error)
	CheckUserExists(ctx context.Context, username string) (bool, error)
}

type tenderRepository struct {
	db *gorm.DB
}

func NewTenderRepository(db *gorm.DB) TenderRepository {
	return &tenderRepository{db: db}
}

func (r *tenderRepository) CreateTender(ctx context.Context, tender *models.Tender) error {
	return r.db.WithContext(ctx).Create(tender).Error
}

func (r *tenderRepository) GetTenderByID(ctx context.Context, id string) (*models.Tender, error) {
	var tender models.Tender
	err := r.db.WithContext(ctx).First(&tender, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tender, nil
}

func (r *tenderRepository) GetTenders(ctx context.Context, serviceTypes []models.TenderServiceType, limit, offset int) ([]*models.Tender, error) {
	var tenders []*models.Tender
	query := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at desc")

	if len(serviceTypes) > 0 {
		query = query.Where("service_type IN ?", serviceTypes)
	}

	err := query.Find(&tenders).Error
	if err != nil {
		return nil, err
	}
	return tenders, nil
}

func (r *tenderRepository) GetTendersByUser(ctx context.Context, username string, limit, offset int) ([]*models.Tender, error) {
	var tenders []*models.Tender

	err := r.db.Table("tenders").
		Joins("JOIN organization_responsible org_resp ON tenders.organization_id = org_resp.organization_id").
		Joins("JOIN employee e ON org_resp.user_id = e.id").
		Where("e.username = ?", username).
		Limit(limit).
		Offset(offset).
		Order("tenders.name").
		Find(&tenders).Error

	if err != nil {
		return nil, err
	}
	return tenders, nil
}

func (r *tenderRepository) UpdateTenderStatus(ctx context.Context, id string, status models.TenderStatus) error {
	result := r.db.Model(&models.Tender{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *tenderRepository) IsUserResponsibleForTender(username, tenderId string) (bool, error) {
	var count int64
	err := r.db.Table("tenders").
		Joins("JOIN organization_responsible org_resp ON tenders.organization_id = org_resp.organization_id").
		Joins("JOIN employee e ON org_resp.user_id = e.id").
		Where("tenders.id = ? AND e.username = ?", tenderId, username).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *tenderRepository) UpdateTender(ctx context.Context, tender *models.Tender) error {
	result := r.db.WithContext(ctx).Model(&models.Tender{}).
		Where("id = ?", tender.ID).
		Updates(map[string]interface{}{
			"name":         tender.Name,
			"description":  tender.Description,
			"service_type": tender.ServiceType,
			"status":       tender.Status,
			"version":      gorm.Expr("version + 1"),
		})

	return result.Error
}

func (r *tenderRepository) DeleteTender(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Tender{}, "id = ?", id).Error
}

func (r *tenderRepository) GetTenderVersions(ctx context.Context, id string) ([]*models.Tender, error) {
	var versions []*models.Tender
	err := r.db.WithContext(ctx).Where("id = ?", id).Order("version asc").Find(&versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

// func (r *tenderRepository) RollbackTenderVersion(ctx context.Context, id string, version int) error {
// 	var versionData models.Tender

// 	// Ищем предыдущую версию по ID тендера и номеру версии
// 	err := r.db.WithContext(ctx).Where("id = ? AND version = ?", id, version).First(&versionData).Error
// 	if err != nil {
// 		return fmt.Errorf("failed to find previous version: %w", err)
// 	}

// 	log.Printf("Previous version data for tender %s: %+v", id, versionData)

// 	newVersion := versionData
// 	newVersion.Version += 1
// 	newVersion.VersionID = "" // GORM сгенерирует новый `version_id`

// 	err = r.db.WithContext(ctx).Create(&newVersion).Error
// 	if err != nil {
// 		return fmt.Errorf("failed to create new version: %w", err)
// 	}

// 	log.Printf("Tender %s successfully rolled back to version %d", id, version)

// 	return nil
// }

func (r *tenderRepository) IsUserResponsibleForOrganization(username, organizationID string) (bool, error) {
	var count int64
	err := r.db.Table("organization_responsible").
		Joins("JOIN employee ON organization_responsible.user_id = employee.id").
		Where("employee.username = ? AND organization_responsible.organization_id = ?", username, organizationID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *tenderRepository) CheckUserExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.Table("employee").Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
