package service

import (
	"fmt"

	"github.com/matthieukhl/align-back/internal/domain"
)

type clientService struct {
	repo domain.ClientRepository
}

// NewClientService creates a new client service
func NewClientService(repo domain.ClientRepository) domain.ClientService {
	return &clientService{
		repo: repo,
	}
}

// GetAll returns all clients
func (s *clientService) GetAll() ([]domain.Client, error) {
	return s.repo.GetAll()
}

// GetByID returns a client by ID
func (s *clientService) GetByID(id string) (*domain.Client, error) {
	return s.repo.GetByID(id)
}

// Create creates a new client
func (s *clientService) Create(input domain.ClientInput) error {
	// Check if email is already used
	existingClient, err := s.repo.GetByEmail(input.Email)
	if err != nil {
		return err
	}

	if existingClient != nil {
		return fmt.Errorf("email %s is already in use", input.Email)
	}

	// Create a new client
	client := &domain.Client{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Phone:          input.Phone,
		Email:          input.Email,
		StreetNumber:   input.StreetNumber,
		StreetName:     input.StreetName,
		City:           input.City,
		ZipCode:        input.ZipCode,
		Country:        input.Country,
		GroupCredits:   input.GroupCredits,
		PrivateCredits: input.PrivateCredits,
	}

	err = s.repo.Create(client)
	if err != nil {
		return err
	}

	return nil
}

// Update updates a client
func (s *clientService) Update(id string, input domain.ClientInput) (*domain.Client, error) {
	// Check if client exists
	existingClient, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existingClient == nil {
		return nil, nil
	}

	// Check if email is already in use by another client
	if existingClient.Email != input.Email {
		clientWithEmail, err := s.repo.GetByEmail(input.Email)
		if err != nil {
			return nil, err
		}

		if clientWithEmail != nil && clientWithEmail.ID != id {
			return nil, fmt.Errorf("email %s is already in use", input.Email)
		}
	}

	// Update client
	client := &domain.Client{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Phone:          input.Phone,
		Email:          input.Email,
		StreetNumber:   input.StreetNumber,
		StreetName:     input.StreetName,
		City:           input.City,
		ZipCode:        input.ZipCode,
		Country:        input.Country,
		GroupCredits:   input.GroupCredits,
		PrivateCredits: input.PrivateCredits,
	}

	err = s.repo.Update(client)
	if err != nil {
		return nil, err
	}

	// Get the updated client to return with all fields
	return s.repo.GetByID(id)
}

// Delete deletes a client
func (s *clientService) Delete(id string) error {
	// Check if client exists
	existingClient, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if existingClient == nil {
		return fmt.Errorf("client with ID %s not found", id)
	}

	return s.repo.Delete(id)
}

// GetByEmail returns a client by email
func (s *clientService) GetByEmail(email string) (*domain.Client, error) {
	return s.repo.GetByEmail(email)
}

// GetLowCredits returns clients with group credits below a threshold
func (s *clientService) GetLowGroupCredits(threshold int) ([]domain.Client, error) {
	return s.repo.GetLowGroupCredits(threshold)
}

// GetLowCredits returns clients with private credits below a threshold
func (s *clientService) GetLowPrivateCredits(threshold int) ([]domain.Client, error) {
	return s.repo.GetLowPrivateCredits(threshold)
}

// UpdateCredits updates a client's group credits
func (s *clientService) UpdateGroupCredits(id string, credits int) error {
	// Check if client exists
	existingClient, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if existingClient == nil {
		return fmt.Errorf("client with id %s not found", id)
	}

	// Update group credits
	existingClient.GroupCredits = credits

	return s.repo.Update(existingClient)
}

// UpdateCredits updates a client's private credits
func (s *clientService) UpdatePrivateCredits(id string, credits int) error {
	// Check if client exists
	existingClient, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if existingClient == nil {
		return fmt.Errorf("client with id %s not found", id)
	}

	// Update group credits
	existingClient.PrivateCredits = credits

	return s.repo.Update(existingClient)
}
