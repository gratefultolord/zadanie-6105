package services

import (
	"context"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/repositories"
)

type BidService struct {
	bidRepo          repositories.BidRepository
	employeeRepo     repositories.EmployeeRepository
	organizationRepo repositories.OrganizationRepository
}

func NewBidService(
	bidRepo repositories.BidRepository,
	employeeRepo repositories.EmployeeRepository,
	organizationRepo repositories.OrganizationRepository,
) *BidService {
	return &BidService{
		bidRepo:          bidRepo,
		employeeRepo:     employeeRepo,
		organizationRepo: organizationRepo,
	}
}

func (s *BidService) IsOrganizationExists(ctx context.Context, organizationID string) (bool, error) {
	return s.organizationRepo.IsOrganizationExists(ctx, organizationID)
}

func (s *BidService) IsEmployeeExists(ctx context.Context, employeeID string) (bool, error) {
	return s.employeeRepo.IsEmployeeExists(ctx, employeeID)
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

func (s *BidService) IsAuthorizedToCreateBid(ctx context.Context, bid *models.Bid, username string, organizationID string) (bool, error) {
	if username != "" && bid.AuthorType == models.AuthorTypeUser {
		return s.isUserAuthorized(ctx, username, bid)
	} else if organizationID != "" && bid.AuthorType == models.AuthorTypeOrganization {
		return s.isOrganizationAuthorized(ctx, organizationID, bid)
	} else {
		return false, nil
	}
}

func (s *BidService) isUserAuthorized(ctx context.Context, username string, bid *models.Bid) (bool, error) {
	userID, err := s.GetAuthorIDByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	if userID == bid.AuthorID {
		return true, nil
	}
	return false, nil
}

func (s *BidService) isOrganizationAuthorized(ctx context.Context, organizationID string, bid *models.Bid) (bool, error) {
	if organizationID == bid.AuthorID {
		return true, nil
	}
	return false, nil
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
	return s.bidRepo.IsUserAuthorizedToViewBids(ctx, tenderID, username)
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
	return s.bidRepo.IsUserAuthorizedForBid(ctx, username, bidID)
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

func (s *BidService) IsUserAuthorizedToDeleteBid(ctx context.Context, username, bidID string) (bool, error) {
	return s.bidRepo.IsUserAuthorizedForBid(ctx, username, bidID)
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

func (s *BidService) AddBidReview(ctx context.Context, review *models.BidReview) error {
	return s.bidRepo.CreateBidReview(ctx, review)
}

func (s *BidService) IsUserAuthorizedToAddReview(ctx context.Context, username string, bidID string) (bool, error) {
	return true, nil
}

func (s *BidService) GetAuthorIDByUsername(ctx context.Context, username string) (string, error) {
	return s.employeeRepo.GetEmployeeIDByUsername(ctx, username)
}
