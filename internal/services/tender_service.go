package services

import (
	"context"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/repositories"
)

type TenderService struct {
	tenderRepo repositories.TenderRepository
}

func NewTenderService(tenderRepo repositories.TenderRepository) *TenderService {
	return &TenderService{tenderRepo: tenderRepo}
}

func (s *TenderService) IsUserAuthorizedToCreateTender(username, organizationID string) (bool, error) {
	return s.tenderRepo.IsUserResponsibleForOrganization(username, organizationID)
}

func (s *TenderService) CreateTender(ctx context.Context, tender *models.Tender) error {
	return s.tenderRepo.CreateTender(ctx, tender)
}

func (s *TenderService) GetTenderByID(ctx context.Context, id string) (*models.Tender, error) {
	return s.tenderRepo.GetTenderByID(ctx, id)
}

func (s *TenderService) GetTendersByUser(ctx context.Context, username string, limit, offset int) ([]*models.Tender, error) {
	return s.tenderRepo.GetTendersByUser(ctx, username, limit, offset)
}

func (s *TenderService) GetTenders(ctx context.Context, serviceTypes []models.TenderServiceType, limit, offset int) ([]*models.Tender, error) {
	return s.tenderRepo.GetTenders(ctx, serviceTypes, limit, offset)
}

func (s *TenderService) UpdateTenderStatus(ctx context.Context, tenderId string, status models.TenderStatus) error {
	return s.tenderRepo.UpdateTenderStatus(ctx, tenderId, status)
}

func (s *TenderService) IsUserAuthorizedToUpdateStatus(username, tenderId string) (bool, error) {
	return s.tenderRepo.IsUserResponsibleForTender(username, tenderId)
}

func (s *TenderService) UpdateTender(ctx context.Context, tenderId string, updates *models.Tender) (*models.Tender, error) {
	existingTender, err := s.tenderRepo.GetTenderByID(ctx, tenderId)
	if err != nil {
		return nil, err
	}

	if updates.Name != "" {
		existingTender.Name = updates.Name
	}
	if updates.Description != "" {
		existingTender.Description = updates.Description
	}
	if updates.ServiceType != "" {
		existingTender.ServiceType = updates.ServiceType
	}

	if err := s.tenderRepo.UpdateTender(ctx, existingTender); err != nil {
		return nil, err
	}

	return existingTender, nil
}

func (s *TenderService) IsUserAuthorizedToEditTender(username, tenderId string) (bool, error) {
	return s.tenderRepo.IsUserResponsibleForTender(username, tenderId)
}

func (s *TenderService) DeleteTender(ctx context.Context, id string) error {
	return s.tenderRepo.DeleteTender(ctx, id)
}

func (s *TenderService) GetTenderVersions(ctx context.Context, id string) ([]*models.Tender, error) {
	return s.tenderRepo.GetTenderVersions(ctx, id)
}

func (s *TenderService) RollbackTenderVersion(ctx context.Context, tenderId string, version int) (*models.Tender, error) {
	if err := s.tenderRepo.RollbackTenderVersion(ctx, tenderId, version); err != nil {
		return nil, err
	}

	// Получаем последнюю версию тендера
	return s.tenderRepo.GetTenderByID(ctx, tenderId)
}

func (s *TenderService) CheckUserExists(ctx context.Context, username string) (bool, error) {
	return s.tenderRepo.CheckUserExists(ctx, username)
}

func (s *TenderService) IsUserAuthorizedToViewStatus(username, tenderId string) (bool, error) {
	return s.tenderRepo.IsUserResponsibleForTender(username, tenderId)
}
