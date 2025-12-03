package dto

import "github.com/google/uuid"

// UpdateLocationRequest represents the request to update driver location
type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

// DriverResponse represents driver information in API responses
type DriverResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Phone     string  `json:"phone"`
	Status    string  `json:"status"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// DriverSearchResponse represents driver information with distance in search results
type DriverSearchResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Phone     string  `json:"phone"`
	Status    string  `json:"status"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  float64 `json:"distance"` // Distance in kilometers
}

// SearchDriversRequest represents the query parameters for searching drivers
type SearchDriversRequest struct {
	Latitude  float64 `form:"lat" binding:"required,min=-90,max=90"`
	Longitude float64 `form:"lng" binding:"required,min=-180,max=180"`
	Radius    float64 `form:"radius" binding:"omitempty,min=0.1,max=100"` // in kilometers, default 5km
}

// ToDriverResponse converts a Driver model to DriverResponse DTO
func ToDriverResponse(id uuid.UUID, name, phone, status string, lat, lng float64) DriverResponse {
	return DriverResponse{
		ID:        id.String(),
		Name:      name,
		Phone:     phone,
		Status:    status,
		Latitude:  lat,
		Longitude: lng,
	}
}
