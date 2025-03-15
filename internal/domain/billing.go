package domain

import "time"

// Billing represents a payment for a package
type Billing struct {
	ID          string    `json:"id" db:"id"`
	ClientID    string    `json:"client_id" db:"client_id"`
	PackageID   string    `json:"package_id" db:"package_id"`
	Amount      int       `json:"amount" db:"amount"`
	Price       float64   `json:"price" db:"price"`
	Credits     int       `json:"credits" db:"credits"`
	PaymentDate time.Time `json:"payment_date" db:"payment_date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Populated from joins
	Client  *Client  `json:"client,omitempty" db:"-"`
	Package *Package `json:"package,omitempty" db:"-"`
}

// BillingWithDetails includes the client and package details
type BillingWithDetails struct {
	Billing *Billing `json:"billing"`
	Client  *Client  `json:"client"`
	Package *Package `json:"package"`
}

// BillingInput is used for creating/updating billings
type BillingInput struct {
	ClientID    string     `json:"client_id" validate:"required,uuid"`
	PackageID   string     `json:"package_id" validate:"required,uuid"`
	Amount      int        `json:"amount" validate:"required,min=1"`
	Price       float64    `json:"price" validate:"required,min=0"`
	PaymentDate *time.Time `json:"payment_date"`
}

// BillingRepository defines methods for billing persistence
type BillingRepository interface {
	GetAll() ([]Billing, error)
	GetByID(id string) (*Billing, error)
	GetByClient(clientID string) ([]Billing, error)
	GetRecent(limit int) ([]Billing, error)
	GetWithDetails(id string) (*BillingWithDetails, error)
	GetAllWithDetails() ([]BillingWithDetails, error)
	Create(billing *Billing) error
	Update(billing *Billing) error
	Delete(id string) error
}

// BillingService defines methods for billing business logic
type BillingService interface {
	GetAll() ([]Billing, error)
	GetByID(id string) (*Billing, error)
	GetByClientID(clientID string) ([]Billing, error)
	GetRecent(limit int) ([]Billing, error)
	GetWithDetails(id string) ([]BillingWithDetails, error)
	GetAllWithDetails() ([]BillingWithDetails, error)
	Create(billing *Billing) error
	Update(billing *Billing) error
	Delete(id string) error
}
