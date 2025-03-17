package service

import (
	"fmt"

	"github.com/matthieukhl/align-back/internal/domain"
)

type billingService struct {
	repo domain.BillingRepository
}

// NewBillingService creates a new billing service
func NewBillingService(repo domain.BillingRepository) domain.BillingService {
	return &billingService{
		repo: repo,
	}
}

// Create creates a new billing.
func (s *billingService) Create(input domain.BillingInput) error {
	// Check if amount is negative or 0
	if input.Amount < 1 {
		return fmt.Errorf("Amount cannot be less than 1: %d", input.Amount)
	}

	// Check if price is negative
	if input.Price < 0.0 {
		return fmt.Errorf("Price cannot negative: %.2f", input.Price)
	}

	// Create new billing
	billing := &domain.Billing{
		ClientID:    input.ClientID,
		PackageID:   input.PackageID,
		Amount:      input.Amount,
		Price:       input.Price,
		Credits:     input.Credits,
		PaymentDate: input.PaymentDate,
	}

	return s.repo.Create(billing)
}

// Delete deletes an existing service.
func (s *billingService) Delete(id string) error {
	// Check if billing ID exists
	existingBilling, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if existingBilling == nil {
		return fmt.Errorf("billing with id %s not found", id)
	}

	return s.repo.Delete(id)
}

// GetAll returns all billings.
func (s *billingService) GetAll() ([]domain.Billing, error) {
	return s.repo.GetAll()
}

// GetAllWithDetails implements domain.BillingService.
func (s *billingService) GetAllWithDetails() ([]domain.BillingWithDetails, error) {
	panic("unimplemented")
}

// GetByClientID retuns billings by client ID.
func (s *billingService) GetByClientID(clientID string) ([]domain.Billing, error) {
	return s.repo.GetByClient(clientID)
}

// GetByID returns a billing by ID.
func (s *billingService) GetByID(id string) (*domain.Billing, error) {
	return s.repo.GetByID(id)
}

// GetRecent returns the most recent billings
func (s *billingService) GetRecent(limit int) ([]domain.Billing, error) {
	return s.repo.GetRecent(limit)
}

// GetWithDetails implements domain.BillingService.
func (s *billingService) GetWithDetails(id string) ([]domain.BillingWithDetails, error) {
	panic("unimplemented")
}

// Update updates an existing billing
func (s *billingService) Update(id string, input domain.BillingInput) (*domain.Billing, error) {
	// Check if billing ID exists
	existingBilling, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existingBilling == nil {
		return nil, fmt.Errorf("billing with ID %s not found", id)
	}

	// Create new billing
	billing := &domain.Billing{
		ClientID:    input.ClientID,
		PackageID:   input.PackageID,
		Amount:      input.Amount,
		Price:       input.Price,
		Credits:     input.Credits,
		PaymentDate: input.PaymentDate,
	}

	err = s.repo.Update(billing)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}
