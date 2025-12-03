package repository

import (
	"trip-service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TripRepository interface {
	CreateTrip(trip *model.Trip) error
	GetTripByID(tripID uuid.UUID) (*model.Trip, error)
	UpdateTrip(trip *model.Trip) error
	GetTripsByPassenger(passengerID uuid.UUID) ([]model.Trip, error)
	GetTripsByDriver(driverID uuid.UUID) ([]model.Trip, error)
}

type tripRepository struct {
	db *gorm.DB
}

func NewTripRepository(db *gorm.DB) TripRepository {
	return &tripRepository{db: db}
}

func (r *tripRepository) CreateTrip(trip *model.Trip) error {
	return r.db.Create(trip).Error
}

func (r *tripRepository) GetTripByID(tripID uuid.UUID) (*model.Trip, error) {
	var trip model.Trip
	err := r.db.First(&trip, "id = ?", tripID).Error
	if err != nil {
		return nil, err
	}
	return &trip, nil
}

func (r *tripRepository) UpdateTrip(trip *model.Trip) error {
	return r.db.Save(trip).Error
}

func (r *tripRepository) GetTripsByPassenger(passengerID uuid.UUID) ([]model.Trip, error) {
	var trips []model.Trip
	err := r.db.Where("passenger_id = ?", passengerID).Order("created_at DESC").Find(&trips).Error
	if err != nil {
		return nil, err
	}
	return trips, nil
}

func (r *tripRepository) GetTripsByDriver(driverID uuid.UUID) ([]model.Trip, error) {
	var trips []model.Trip
	err := r.db.Where("driver_id = ?", driverID).Order("created_at DESC").Find(&trips).Error
	if err != nil {
		return nil, err
	}
	return trips, nil
}
