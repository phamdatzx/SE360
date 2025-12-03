package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Driver status constants
const (
	StatusAvailable = "available"
	StatusBusy      = "busy"
	StatusOffline   = "offline"
)

type Driver struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	Name      string    `json:"name" gorm:"column:name"`
	Phone     string    `json:"phone" gorm:"column:phone;uniqueIndex"`
	Status    string    `json:"status" gorm:"column:status;default:offline"` // available, busy, offline
	Latitude  float64   `json:"latitude" gorm:"column:latitude"`
	Longitude float64   `json:"longitude" gorm:"column:longitude"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// Hook to generate UUID before creating record
func (d *Driver) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = uuid.New()
	return
}

// TableName specifies the table name for Driver model
func (Driver) TableName() string {
	return "drivers"
}
