package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Trip status constants
const (
	StatusSearchingDriver = "searching_driver"
	StatusAccepted        = "accepted"
	StatusInProgress      = "in_progress"
	StatusCompleted       = "completed"
	StatusCancelled       = "cancelled"
)

type Trip struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;column:id"`
	PassengerID     uuid.UUID  `gorm:"type:uuid;column:passenger_id;index"`
	DriverID        *uuid.UUID `gorm:"type:uuid;column:driver_id;index"`
	PickupLatitude  float64    `json:"pickup_latitude" gorm:"column:pickup_latitude"`
	PickupLongitude float64    `json:"pickup_longitude" gorm:"column:pickup_longitude"`
	DropoffLatitude float64    `json:"dropoff_latitude" gorm:"column:dropoff_latitude"`
	DropoffLongitude float64   `json:"dropoff_longitude" gorm:"column:dropoff_longitude"`
	Status          string     `json:"status" gorm:"column:status;index"`
	Fare            *float64   `json:"fare" gorm:"column:fare"`
	CancelReason    *string    `json:"cancel_reason" gorm:"column:cancel_reason"`
	CreatedAt       time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"column:updated_at"`
}

// Hook to generate UUID before creating record
func (t *Trip) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}

// TableName specifies the table name for Trip model
func (Trip) TableName() string {
	return "trips"
}
