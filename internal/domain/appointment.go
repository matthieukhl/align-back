package domain

import (
	"time"
)

// Appointment represents a client booking for a scheduled class
type Appointment struct {
	ID         string    `json:"id" db:"id"`
	ScheduleID string    `json:"schedule_id" db:"schedule_id"`
	ClientID   string    `json:"client_id" db:"client_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`

	// Populated from joins
	Client   *Client   `json:"client,omitempty" db:"-"`
	Schedule *Schedule `json:"schedule,omitempty" db:"-"`
}

// AppointmentWithDetails includes the schedule and client details
type AppointmentWithDetails struct {
	Appointment *Appointment         `json:"appointment"`
	Client      *Client              `json:"client"`
	Schedule    *ScheduleWithDetails `json:"schedule"`
}

// AppointmentInput is used for creating/updating appointments
type AppointmentInput struct {
	ScheduleID string `json:"schedule_id" validate:"required,uuid"`
	ClientID   string `json:"client_id" validate:"required,uuid"`
}

// ApointmentRepository defines methods for appointment persistence
type AppointmentRepository interface {
	GetAll() ([]Appointment, error)
	GetByID(id string) (*Appointment, error)
	GetByClient(clientID string) ([]Appointment, error)
	GetBySchedule(scheduleID string) ([]Appointment, error)
	GetByClientAndSchedule(clientID, scheduleID string) (*Appointment, error)
	GetUpcomingByClient(clientID string) ([]Appointment, error)
	GetWithDetails(id string) (*AppointmentWithDetails, error)
	Create(appointment *Appointment) error
	Update(appointment *Appointment) error
	Delete(id string) error
	CountBySchedule(scheduleID string) (int, error)
}

// AppointmentService defines methods for appointment business logic
type AppointmentService interface {
	GetAll() ([]Appointment, error)
	GetByID(id string) (*Appointment, error)
	GetByClientID(clientID string) ([]Appointment, error)
	GetByScheduleID(scheduleID string) ([]Appointment, error)
	GetUpcomingByClient(clientID string) ([]Appointment, error)
	GetWithDetails(id string) (*AppointmentWithDetails, error)
	Create(appointment *Appointment) error
	Update(appointment *Appointment) error
	Delete(id string) error
}
