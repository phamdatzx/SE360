package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateTripRequest represents the request to create a new trip
type CreateTripRequest struct {
	PassengerID      string  `json:"passenger_id" binding:"required,uuid"`
	PickupLatitude   float64 `json:"pickup_lat" binding:"required,min=-90,max=90"`
	PickupLongitude  float64 `json:"pickup_lng" binding:"required,min=-180,max=180"`
	DropoffLatitude  float64 `json:"dropoff_lat" binding:"required,min=-90,max=90"`
	DropoffLongitude float64 `json:"dropoff_lng" binding:"required,min=-180,max=180"`
}

// TripResponse represents trip information in API responses
type TripResponse struct {
	ID               string     `json:"id"`
	PassengerID      string     `json:"passenger_id"`
	DriverID         *string    `json:"driver_id"`
	PickupLatitude   float64    `json:"pickup_latitude"`
	PickupLongitude  float64    `json:"pickup_longitude"`
	DropoffLatitude  float64    `json:"dropoff_latitude"`
	DropoffLongitude float64    `json:"dropoff_longitude"`
	Status           string     `json:"status"`
	Fare             *float64   `json:"fare"`
	CancelReason     *string    `json:"cancel_reason"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// CancelTripRequest represents the request to cancel a trip
type CancelTripRequest struct {
	Reason string `json:"reason" binding:"omitempty,max=500"`
}

// UpdateTripStatusRequest represents the request to update trip status
type UpdateTripStatusRequest struct {
	Status   string  `json:"status" binding:"required,oneof=searching_driver accepted in_progress completed cancelled"`
	DriverID *string `json:"driver_id" binding:"omitempty,uuid"`
	Fare     *float64 `json:"fare" binding:"omitempty,min=0"`
}

// ToTripResponse converts a Trip model to TripResponse DTO
func ToTripResponse(
	id, passengerID uuid.UUID,
	driverID *uuid.UUID,
	pickupLat, pickupLng, dropoffLat, dropoffLng float64,
	status string,
	fare *float64,
	cancelReason *string,
	createdAt, updatedAt time.Time,
) TripResponse {
	var driverIDStr *string
	if driverID != nil {
		str := driverID.String()
		driverIDStr = &str
	}

	return TripResponse{
		ID:               id.String(),
		PassengerID:      passengerID.String(),
		DriverID:         driverIDStr,
		PickupLatitude:   pickupLat,
		PickupLongitude:  pickupLng,
		DropoffLatitude:  dropoffLat,
		DropoffLongitude: dropoffLng,
		Status:           status,
		Fare:             fare,
		CancelReason:     cancelReason,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}
