package services

import (
	"context"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/repositories"
)

type BidService struct {
	bidRepo repositories.BidRepository
}

func NewBidService(bidRepo repositories.BidRepository) *BidService {
	return &BidService{bidRepo: bidRepo}
}

func (s *BidService) CreateBid(ctx context.Context, bid *models.Bid) error {
	return s.bidRepo.CreateBid(ctx, bid)
}

func (s *BidService) IsTenderExists(ctx context.Context, tenderID string) (bool, error) {
	return s.bidRepo.IsTenderExists(ctx, tenderID)
}

func (s *BidService) IsUserAuthorizedToCreateBid(ctx context.Context, username, tenderID string) (bool, error) {
	return s.bidRepo.IsUserResponsibleForTender(ctx, username, tenderID)
}

func (s *BidService) GetBid(ctx context.Context, id string) (*models.Bid, error) {
	return s.bidRepo.GetBidByID(ctx, id)
}

func (s *BidService) GetUserBids(ctx context.Context, username string, limit, offset int) ([]*models.Bid, error) {
	return s.bidRepo.GetBidsByUser(ctx, username, limit, offset)
}

func (s *BidService) GetBidsForTender(ctx context.Context, tenderID string, limit, offset int) ([]*models.Bid, error) {
	return s.bidRepo.GetBidsForTender(ctx, tenderID, limit, offset)
}

func (s *BidService) IsUserAuthorizedToViewBids(ctx context.Context, username, tenderID string) (bool, error) {
	return s.bidRepo.IsUserResponsibleForTender(ctx, username, tenderID)
}

func (s *BidService) GetAllBids(ctx context.Context) ([]*models.Bid, error) {
	return s.bidRepo.GetAllBids(ctx)
}

func (s *BidService) GetBidStatus(ctx context.Context, bidID string) (string, error) {
	return s.bidRepo.GetBidStatus(ctx, bidID)
}

func (s *BidService) IsUserAuthorizedToViewBid(ctx context.Context, username, bidID string) (bool, error) {
	return s.bidRepo.IsUserAuthorizedForBid(ctx, username, bidID)
}

func (s *BidService) UpdateBidStatus(ctx context.Context, bidID string, status string) error {
	return s.bidRepo.UpdateBidStatus(ctx, bidID, status)
}

func (s *BidService) IsUserAuthorizedToChangeStatus(ctx context.Context, username, bidID string) (bool, error) {
	return s.bidRepo.IsUserResponsibleForBid(ctx, username, bidID)
}

func (s *BidService) UpdateBid(ctx context.Context, bidID string, updatedBid *models.Bid) (*models.Bid, error) {
	existingBid, err := s.bidRepo.GetBidByID(ctx, bidID)
	if err != nil {
		return nil, err
	}

	if updatedBid.Name != "" {
		existingBid.Name = updatedBid.Name
	}
	if updatedBid.Description != "" {
		existingBid.Description = updatedBid.Description
	}

	existingBid.Version++

	err = s.bidRepo.UpdateBid(ctx, existingBid)
	if err != nil {
		return nil, err
	}

	return existingBid, nil
}

func (s *BidService) IsUserAuthorizedToEditBid(ctx context.Context, username, bidID string) (bool, error) {
	return s.bidRepo.IsUserResponsibleForBid(ctx, username, bidID)
}

func (s *BidService) SubmitBidDecision(ctx context.Context, bidID string, decision string) error {
	bid, err := s.bidRepo.GetBidByID(ctx, bidID)
	if err != nil {
		return err
	}

	if decision == "Approved" {
		bid.Status = models.BidStatusPublished
	} else if decision == "Rejected" {
		bid.Status = models.BidStatusCanceled
	}

	err = s.bidRepo.UpdateBid(ctx, bid)
	if err != nil {
		return err
	}

	return nil
}

func (s *BidService) IsUserAuthorizedToSubmitDecision(ctx context.Context, username, bidID string) (bool, error) {
	return s.bidRepo.IsUserResponsibleForBid(ctx, username, bidID)
}

func (s *BidService) DeleteBid(ctx context.Context, id string) error {
	return s.bidRepo.DeleteBid(ctx, id)
}

func (s *BidService) SubmitBidFeedback(ctx context.Context, bidID string, feedback string) error {
	bid, err := s.bidRepo.GetBidByID(ctx, bidID)
	if err != nil {
		return err
	}

	bid.Feedback = feedback

	err = s.bidRepo.UpdateBid(ctx, bid)
	if err != nil {
		return err
	}

	return nil
}

func (s *BidService) IsUserAuthorizedToSubmitFeedback(ctx context.Context, username, bidID string) (bool, error) {
	return s.bidRepo.IsUserResponsibleForBid(ctx, username, bidID)
}

func (s *BidService) RollbackBidVersion(ctx context.Context, bidID string, version int) error {
	oldBid, err := s.bidRepo.GetBidByVersion(ctx, bidID, version)
	if err != nil {
		return err
	}

	oldBid.Version++

	err = s.bidRepo.UpdateBid(ctx, oldBid)
	if err != nil {
		return err
	}

	return nil
}

func (s *BidService) IsUserAuthorizedToRollback(ctx context.Context, username, bidID string) (bool, error) {
	return s.bidRepo.IsUserResponsibleForBid(ctx, username, bidID)
}

func (s *BidService) GetBidReviews(ctx context.Context, tenderID string, authorUsername string, limit, offset int) ([]*models.BidReview, error) {
	return s.bidRepo.GetBidReviews(ctx, tenderID, authorUsername, limit, offset)
}

func (s *BidService) IsUserAuthorizedToViewReviews(ctx context.Context, username, tenderID string) (bool, error) {
	return s.bidRepo.IsUserResponsibleForTender(ctx, username, tenderID)
}
