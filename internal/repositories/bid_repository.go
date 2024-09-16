package repositories

import (
	"context"
	"zadanie-6105/internal/models"

	"gorm.io/gorm"
)

type BidRepository interface {
	CreateBid(ctx context.Context, bid *models.Bid) error
	IsTenderExists(ctx context.Context, tenderID string) (bool, error)
	IsUserResponsibleForTender(ctx context.Context, username, tenderID string) (bool, error)
	GetBidByID(ctx context.Context, id string) (*models.Bid, error)
	GetBidsByUser(ctx context.Context, username string, limit, offset int) ([]*models.Bid, error)
	GetBidsForTender(ctx context.Context, tenderID string, limit, offset int) ([]*models.Bid, error)
	GetAllBids(ctx context.Context) ([]*models.Bid, error)
	GetBidStatus(ctx context.Context, bidID string) (string, error)
	IsUserAuthorizedForBid(ctx context.Context, username, bidID string) (bool, error)
	UpdateBidStatus(ctx context.Context, bidID string, status string) error
	IsUserResponsibleForBid(ctx context.Context, username, bidID string) (bool, error)
	UpdateBid(ctx context.Context, bid *models.Bid) error
	DeleteBid(ctx context.Context, id string) error
	GetBidByVersion(ctx context.Context, bidID string, version int) (*models.Bid, error)
	GetBidReviews(ctx context.Context, tenderID string, authorUsername string, limit, offset int) ([]*models.BidReview, error)
}

type bidRepository struct {
	db *gorm.DB
}

func NewBidRepository(db *gorm.DB) BidRepository {
	return &bidRepository{db: db}
}

func (r *bidRepository) CreateBid(ctx context.Context, bid *models.Bid) error {
	return r.db.WithContext(ctx).Create(bid).Error
}

func (r *bidRepository) IsTenderExists(ctx context.Context, tenderID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Tender{}).
		Where("id = ?", tenderID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *bidRepository) IsUserResponsibleForTender(ctx context.Context, username, tenderID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("tenders").
		Joins("JOIN organization_responsible or ON tenders.organization_id = or.organization_id").
		Joins("JOIN employee e ON or.user_id = e.id").
		Where("tenders.id = ? AND e.username = ?", tenderID, username).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *bidRepository) GetBidByID(ctx context.Context, id string) (*models.Bid, error) {
	var bid models.Bid
	err := r.db.WithContext(ctx).First(&bid, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &bid, nil
}

func (r *bidRepository) GetBidsByUser(ctx context.Context, username string, limit, offset int) ([]*models.Bid, error) {
	var bids []*models.Bid

	err := r.db.WithContext(ctx).
		Table("bids").
		Joins("JOIN employee ON bids.author_id = employee.id").
		Where("employee.username = ?", username).
		Limit(limit).
		Offset(offset).
		Order("bids.name").
		Find(&bids).Error

	if err != nil {
		return nil, err
	}

	return bids, nil
}

func (r *bidRepository) GetBidsForTender(ctx context.Context, tenderID string, limit, offset int) ([]*models.Bid, error) {
	var bids []*models.Bid

	err := r.db.WithContext(ctx).
		Where("tender_id = ?", tenderID).
		Limit(limit).
		Offset(offset).
		Order("name").
		Find(&bids).Error

	if err != nil {
		return nil, err
	}

	return bids, nil
}

func (r *bidRepository) GetAllBids(ctx context.Context) ([]*models.Bid, error) {
	var bids []*models.Bid
	err := r.db.WithContext(ctx).Find(&bids).Error
	if err != nil {
		return nil, err
	}
	return bids, nil
}

func (r *bidRepository) GetBidStatus(ctx context.Context, bidID string) (string, error) {
	var status string

	err := r.db.WithContext(ctx).
		Table("bids").
		Select("status").
		Where("id = ?", bidID).
		Scan(&status).Error

	if err != nil {
		return "", err
	}

	return status, nil
}

func (r *bidRepository) IsUserAuthorizedForBid(ctx context.Context, username, bidID string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Table("bids").
		Joins("JOIN employee e ON bids.author_id = e.id").
		Where("bids.id = ? AND e.username = ?", bidID, username).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *bidRepository) UpdateBidStatus(ctx context.Context, bidID string, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.Bid{}).
		Where("id = ?", bidID).
		Update("status", status).Error
}

func (r *bidRepository) IsUserResponsibleForBid(ctx context.Context, username, bidID string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Table("bids").
		Joins("JOIN employee e ON bids.author_id = e.id").
		Where("bids.id = ? AND e.username = ?", bidID, username).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *bidRepository) UpdateBid(ctx context.Context, bid *models.Bid) error {
	return r.db.WithContext(ctx).
		Model(&models.Bid{}).
		Where("id = ?", bid.ID).
		Updates(map[string]interface{}{
			"name":        bid.Name,
			"description": bid.Description,
			"version":     bid.Version,
		}).Error
}

func (r *bidRepository) DeleteBid(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Bid{}, "id = ?", id).Error
}

func (r *bidRepository) GetBidByVersion(ctx context.Context, bidID string, version int) (*models.Bid, error) {
	var bid models.Bid

	err := r.db.WithContext(ctx).
		Where("id = ? AND version = ?", bidID, version).
		First(&bid).Error
	if err != nil {
		return nil, err
	}

	return &bid, nil
}

func (r *bidRepository) GetBidReviews(ctx context.Context, tenderID string, authorUsername string, limit, offset int) ([]*models.BidReview, error) {
	var reviews []*models.BidReview

	err := r.db.WithContext(ctx).
		Table("bid_reviews").
		Joins("JOIN bids ON bid_reviews.bid_id = bids.id").
		Joins("JOIN employee e ON bids.author_id = e.id").
		Where("bids.tender_id = ? AND e.username = ?", tenderID, authorUsername).
		Limit(limit).
		Offset(offset).
		Find(&reviews).Error

	if err != nil {
		return nil, err
	}

	return reviews, nil
}
