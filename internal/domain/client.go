package domain

import "time"

// Client represents a pilates client
type Client struct {
	ID             string    `json:"id" db:"id"`
	FirstName      string    `json:"firstname" db:"firstname"`
	LastName       string    `json:"lastname" db:"lastname"`
	Phone          string    `json:"phone" db:"phone"`
	Email          string    `json:"email" db:"email"`
	StreetNumber   string    `json:"street_number" db:"street_number"`
	StreetName     string    `json:"street_name" db:"street_name"`
	City           string    `json:"city" db:"city"`
	ZipCode        string    `json:"zip_code" db:"zip_code"`
	Country        string    `json:"country" db:"country"`
	GroupCredits   int       `json:"group_credits" db:"group_credits"`
	PrivateCredits int       `json:"private_credits" db:"private_credits"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// ClientInput is used for creating/updating clients
type ClientInput struct {
	FirstName      string `json:"firstname" validate:"required"`
	LastName       string `json:"lastname" validate:"required"`
	Phone          string `json:"phone"`
	Email          string `json:"email" validate:"required,email"`
	StreetNumber   string `json:"street_number"`
	StreetName     string `json:"street_name"`
	City           string `json:"city"`
	ZipCode        string `json:"zip_code"`
	Country        string `json:"country"`
	GroupCredits   int    `json:"group_credits" db:"group_credits"`
	PrivateCredits int    `json:"private_credits" db:"private_credits"`
}

// ClientRepository defines methods for client persistence
type ClientRepository interface {
	GetAll() ([]Client, error)
	GetByID(id string) (*Client, error)
	GetByEmail(email string) (*Client, error)
	GetLowGroupCredits(threshold int) ([]Client, error)
	GetLowPrivateCredits(threshold int) ([]Client, error)
	Create(client *Client) error
	Update(client *Client) error
	Delete(id string) error
}

// ClientService defines business logic for clients
type ClientService interface {
	GetAll() ([]Client, error)
	GetByID(id string) (*Client, error)
	Create(input ClientInput) error
	Update(id string, input ClientInput) (*Client, error)
	Delete(id string) error
	GetByEmail(email string) (*Client, error)
	GetLowGroupCredits(threshold int) ([]Client, error)
	GetLowPrivateCredits(threshold int) ([]Client, error)
	UpdateGroupCredits(id string, groupCredits int) error
	UpdatePrivateCredits(id string, privateCredits int) error
}
