package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/matthieukhl/align-back/internal/domain"
	"github.com/rs/zerolog/log"
)

type appointmentRepository struct {
	db *sqlx.DB
}

// CountBySchedule return the count of appoitments for a given schedule ID.
func (r *appointmentRepository) CountBySchedule(scheduleID string) (int, error) {
	var count int

	query := `
	SELECT COUNT(1) FROM appointments WHERE schedule_id = ?
	`

	err := r.db.Get(&count, query, scheduleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		log.Error().Err(err).Str("scheduleID", scheduleID).Msg("failed to retrieve appointment count by schedule")
		return 0, fmt.Errorf("failed to retrieve appointment count by schedule: %w", err)
	}

	return count, nil
}

// Create creates a new appointment.
func (r *appointmentRepository) Create(appointment *domain.Appointment) error {
	query := `
	INSERT INTO
		appointments (
			schedule_id
			, client_id
		)
	VALUES (?, ?)
	`

	_, err := r.db.Exec(query, appointment.ScheduleID, appointment.ClientID)
	if err != nil {
		log.Error().Err(err).Interface("appointment", appointment).Msg("failed to create appointment")
		return fmt.Errorf("failed to create appointment: %w", err)
	}

	return nil
}

// Delete deletes an appointment.
func (r *appointmentRepository) Delete(id string) error {
	query := `
	DELETE FROM
		appointments
	WHERE 
		id = ?
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to delete appointment")
		return fmt.Errorf("failed to delete appointment: %w", err)
	}

	return nil
}

// GetAll returns all appointments.
func (r *appointmentRepository) GetAll() ([]domain.Appointment, error) {
	var appointments []domain.Appointment

	query := `
	SELECT
		id
		, schedule_id
		, client_id
	FROM
		appointments
	`

	err := r.db.Select(&appointments, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all appointments")
		return nil, fmt.Errorf("failed to get all appointments: %w", err)
	}

	return appointments, nil
}

// GetByClient returns a list of appointments for a given client ID.
func (r *appointmentRepository) GetByClient(clientID string) ([]domain.Appointment, error) {
	var appointments []domain.Appointment

	query := `
	SELECT
		id
		, schedule_id
		, client_id
	FROM
		appointments
	WHERE 
		client_id = ?
	`

	err := r.db.Select(&appointments, query, clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Str("clientID", clientID).Msg("failed to retrieve appointments by client ID")
		return nil, fmt.Errorf("failed to retrieve appointments by ID: %w", err)
	}

	return appointments, nil
}

// GetByClientAndSchedule returns an appointment given a client ID and a schedule ID.
func (r *appointmentRepository) GetByClientAndSchedule(clientID string, scheduleID string) (*domain.Appointment, error) {
	var appointment *domain.Appointment

	query := `
		SELECT 
			id
			, schedule_id
			, client_id
		FROM 
			appointments
		WHERE
			schedule_id = ?
			AND client_id = ?
	`

	err := r.db.Get(&appointment, query, scheduleID, clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Str("clientID", clientID).Str("scheduleId", scheduleID).Msg("failed to get appointment by schedule and client ID")
		return nil, fmt.Errorf("failed to get appointment by client ID and schedule ID: %w", err)
	}

	return appointment, nil
}

// GetByID returns an appointment by its ID.
func (r *appointmentRepository) GetByID(id string) (*domain.Appointment, error) {
	var appointment *domain.Appointment

	query := `
	SELECT
		id
		, schedule_id
		, client_id
	FROM
		appointments
	WHERE 
		id = ?
	`

	err := r.db.Get(&appointment, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Str("id", id).Msg("failed to get appointment by ID")
		return nil, fmt.Errorf("failed to get appointment by ID: %w", err)
	}

	return appointment, nil
}

// GetBySchedule implements domain.AppointmentRepository.
func (r *appointmentRepository) GetBySchedule(scheduleID string) ([]domain.Appointment, error) {
	panic("unimplemented")
}

// GetUpcomingByClient implements domain.AppointmentRepository.
func (r *appointmentRepository) GetUpcomingByClient(clientID string) ([]domain.Appointment, error) {
	panic("unimplemented")
}

// GetWithDetails implements domain.AppointmentRepository.
func (r *appointmentRepository) GetWithDetails(id string) (*domain.AppointmentWithDetails, error) {
	panic("unimplemented")
}

// Update implements domain.AppointmentRepository.
func (r *appointmentRepository) Update(appointment *domain.Appointment) error {
	panic("unimplemented")
}

// NewAppointmentRepository creates a new appointment repository
func NewAppointmentRepository(db *sqlx.DB) domain.AppointmentRepository {
	return &appointmentRepository{
		db: db,
	}
}
