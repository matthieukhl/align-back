package service

import (
	"fmt"

	"github.com/matthieukhl/align-back/internal/domain"
)

type packageService struct {
	repo domain.PackageRepository
}

// NewPackageService creates a new package service
func NewPackageService(repo domain.PackageRepository) domain.PackageService {
	return &packageService{
		repo: repo,
	}
}

// Create creates a new package.
func (s *packageService) Create(input domain.PackageInput) error {
	// Check if package already exists
	existingPackage, err := s.repo.GetByName(input.Name)
	if err != nil {
		return err
	}

	if existingPackage != nil {
		return fmt.Errorf("package %s already exists", input.Name)
	}

	// Create a new package
	pkg := &domain.Package{
		Name:             input.Name,
		NumberOfSessions: input.NumberOfSessions,
		Type:             input.Type,
		Price:            input.Price,
	}

	err = s.repo.Create(pkg)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes an existing package.
func (s *packageService) Delete(id string) error {
	// Check if package exists
	existingPackage, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if existingPackage == nil {
		return fmt.Errorf("package with ID %s not found", id)
	}

	return s.repo.Delete(id)
}

// GetAll returns all packages.
func (s *packageService) GetAll() ([]domain.Package, error) {
	return s.repo.GetAll()
}

// GetByID returns a package by ID.
func (s *packageService) GetByID(id string) (*domain.Package, error) {
	return s.repo.GetByID(id)
}

// GetByType returns a package by type.
func (s *packageService) GetByType(pkgType domain.PackageType) ([]domain.Package, error) {
	return s.repo.GetByType(pkgType)
}

// GetByName returns a package by name.
func (s *packageService) GetByName(name string) (*domain.Package, error) {
	return s.repo.GetByName(name)
}

// Update implements domain.PackageService.
func (s *packageService) Update(id string, input domain.PackageInput) (*domain.Package, error) {
	// Check if package exists
	existingPackage, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existingPackage == nil {
		return nil, fmt.Errorf("package with id %s not found", id)
	}

	// Check if name is already in use by another package
	if existingPackage.Name != input.Name {
		packageWithName, err := s.repo.GetByName(input.Name)
		if err != nil {
			return nil, err
		}

		if packageWithName != nil && packageWithName.ID != id {
			return nil, fmt.Errorf("package name %s already exists")
		}
	}

	// Update package
	pkg := &domain.Package{
		Name:             input.Name,
		NumberOfSessions: input.NumberOfSessions,
		Type:             input.Type,
		Price:            input.Price,
	}

	err = s.repo.Update(pkg)
	if err != nil {
		return nil, err
	}

	// Get the updated package to return with all fields
	return s.repo.GetByID(id)
}
