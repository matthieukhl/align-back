package domain

import "time"

// Schedule represents a scheduled class
type Schedule struct {
	ID            string    `json:"id" db:"id"`
	ClassID       string    `json:"class_id" db:"class_id"`
	Capacity      int       `json:"capacity" db:"capacity"`
	ClassDatetime time.Time `json:"class_datetime" db:"class_datetime"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`

	// Populated from joins
	Class       *Class `json:"class,omitempty" db:"-"`
	BookedCount int    `json:"booked_count,omitempty" db:"-"`
}

// ScheduleWithDetails includes the class details and booking count
type ScheduleWithDetails struct {
	Schedule       *Schedule `json:"schedule"`
	Class          *Class    `json:"class"`
	BookedCount    int       `json:"booked_count"`
	AvailableSlots int       `json:"available_slots"`
}

// ScheduleInput is used for creating/updating schedules
type ScheduleInput struct {
	ClassID       string    `json:"class_id" validate:"required,uuid"`
	Capacity      int       `json:"capacity" validate:"required,min=1"`
	ClassDatetime time.Time `json:"class_datetime" validate:"required"`
}

// ScheduleRepository defines methods for schedule persistence
type ScheduleRepository interface {
	GetAll() ([]Schedule, error)
	GetByID(id string) (*Schedule, error)
	GetByDate(date time.Time) ([]Schedule, error)
	GetByDateRange(startDate, endDate time.Time) ([]Schedule, error)
	GetByClass(classID string) ([]Schedule, error)
	GetUpcoming(limit int) ([]Schedule, error)
	GetWithDetails(id string) (*ScheduleWithDetails, error)
	GetAllWithDetails() ([]ScheduleWithDetails, error)
	Create(schedule *Schedule) error
	Update(schedule *Schedule) error
	Delete(id string) error
}

// ScheduleService defines methods for schedule business logic
type ScheduleService interface {
	GetAll() ([]Schedule, error)
	GetByID(id string) (*Schedule, error)
	GetByDate(date time.Time) ([]Schedule, error)
	GetByWeek(date time.Time) ([]Schedule, error)
	GetByClass(classID string) ([]Schedule, error)
	GetUpcoming(limit int) ([]Schedule, error)
	GetWithDetails(id string) (*ScheduleWithDetails, error)
	GetAllWithDetails() ([]ScheduleWithDetails, error)
	Create(input ScheduleInput) (*Schedule, error)
	Update(id string, input ScheduleInput) (*Schedule, error)
	Delete(id string) error
}
