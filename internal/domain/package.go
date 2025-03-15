package domain

import "time"

// PackageType represents the type of package
type PackageType string

const (
	GroupPackage   PackageType = "GROUP"
	PrivatePackage PackageType = "PRIVATE"
)

// Package represents a Pilates session Package
type Package struct {
	ID               string      `json:"id" db:"id"`
	Name             string      `json:"name" db:"name"`
	NumberOfSessions int         `json:"number_of_sessions" db:"number_of_sessions"`
	Type             PackageType `json:"type" db:"type"`
	Price            float64     `json:"price" db:"price"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" db:"updated_at"`
}

// PackageInput is used for creating/updating packages
type PackageInput struct {
	Name             string      `json:"name" validate:"required"`
	NumberOfSessions int         `json:"number_of_sessions" validate:"required,min=1"`
	Type             PackageType `json:"type" validate:"required,oneof=GROUP PRIVATE"`
	Price            float64     `json:"price" validate:"required,min=0"`
}

// PackageRepository defines methods for package persistence
type PackageRepository interface {
	GetAll() ([]Package, error)
	GetByID(id string) (*Package, error)
	GetByName(name string) (*Package, error)
	GetByType(pkgType PackageType) ([]Package, error)
	Create(pkg *Package) error
	Update(pkg *Package) error
	Delete(id string) error
}

// PackageService defines business logic for packages
type PackageService interface {
	GetAll() ([]Package, error)
	GetByID(id string) (*Package, error)
	GetByName(name string) (*Package, error)
	GetByType(pkgType PackageType) ([]Package, error)
	Create(input PackageInput) error
	Update(id string, input PackageInput) (*Package, error)
	Delete(id string) error
}
