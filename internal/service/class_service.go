package service

import (
	"fmt"

	"github.com/matthieukhl/align-back/internal/domain"
)

type classService struct {
	repo domain.ClassRepository
}

// NewClassService creates a new class service
func NewClassService(repo domain.ClassRepository) domain.ClassService {
	return &classService{
		repo: repo,
	}
}

// GetAll returns all classes
func (s *classService) GetAll() ([]domain.Class, error) {
	return s.repo.GetAll()
}

// GetByID returns a class by ID
func (s *classService) GetByID(id string) (*domain.Class, error) {
	return s.repo.GetByID(id)
}

// GetByLocation returns classes by location
func (s *classService) GetByLocation(location domain.Location) ([]domain.Class, error) {
	return s.repo.GetByLocation(location)
}

// GetByType returns classes by type
func (s *classService) GetByType(classType domain.ClassType) ([]domain.Class, error) {
	return s.repo.GetByType(classType)
}

// Create creates a class
func (s *classService) Create(input domain.ClassInput) error {
	// Check if class name is already used
	existingClass, err := s.repo.GetByName(input.Name)
	if err != nil {
		return err
	}

	if existingClass != nil {
		return fmt.Errorf("class name %s is already in use", input.Name)
	}

	// Create a new class
	class := &domain.Class{
		Name:      input.Name,
		Location:  input.Location,
		Type:      input.Type,
		Equipment: input.Equipment,
	}

	err = s.repo.Create(class)
	if err != nil {
		return err
	}

	return nil
}

// Updates update an existing class
func (s *classService) Update(id string, input domain.ClassInput) (*domain.Class, error) {
	// Check if class ID exists
	existingClass, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existingClass == nil {
		return nil, err
	}

	// Check if name is already in use by another client
	if existingClass.Name != input.Name {
		classWithName, err := s.repo.GetByName(input.Name)
		if err != nil {
			return nil, err
		}

		if classWithName != nil && classWithName.ID != id {
			return nil, fmt.Errorf("class name %s is already in use", input.Name)
		}
	}

	// Update client
	class := &domain.Class{
		Name:      input.Name,
		Location:  input.Location,
		Type:      input.Type,
		Equipment: input.Equipment,
	}

	err = s.repo.Update(class)
	if err != nil {
		return nil, err
	}

	// Get the updated class to return with all fields
	return s.repo.GetByID(id)

}

// Delete deletes a class
func (s *classService) Delete(id string) error {
	// Check if class exists
	existingClass, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if existingClass == nil {
		return fmt.Errorf("class with ID %s not found", id)
	}

	return s.repo.Delete(id)
}
