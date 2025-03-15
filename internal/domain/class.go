package domain

import (
	"time"
)

// Location represents the location where a class is held
type Location string

const (
	Clairvivre Location = "CLAIRVIVRE"
	Cubjac     Location = "CUBJAC"
)

// ClassType represents the type of class
type ClassType string

const (
	GroupClass   ClassType = "GROUP"
	PrivateClass ClassType = "PRIVATE"
)

// Class represents a pilates class
type Class struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Location  Location  `json:"location" db:"location"`
	Type      ClassType `json:"type" db:"type"`
	Equipment string    `json:"equipment" db:"equipment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ClassInput is used for creating/updating classes
type ClassInput struct {
	Name      string    `json:"name" validate:"required"`
	Location  Location  `json:"location" validate:"required,oneof=CLAIRVIVRE CUBJAC"`
	Type      ClassType `json:"type" validate:"required,oneof=GROUP PRIVATE"`
	Equipment string    `json:"equipment"`
}

// ClassRepository defines methods for class persistence
type ClassRepository interface {
	GetAll() ([]Class, error)
	GetByID(id string) (*Class, error)
	GetByType(classType ClassType) ([]Class, error)
	GetByLocation(location Location) ([]Class, error)
	Create(class *Class) error
	Update(class *Class) error
	Delete(id string) error
}

// ClassService defines methods for class business logic
type ClassService interface {
	GetAll() ([]Class, error)
	GetByID(id string) (*Class, error)
	GetByType(classType ClassType) ([]Class, error)
	GetByLocation(location Location) ([]Class, error)
	Create(class *Class) error
	Update(class *Class) error
	Delete(id string) error
}
