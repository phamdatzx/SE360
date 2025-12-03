package service

import (
	"driver-service/dto"
	"driver-service/repository"

	"github.com/google/uuid"
)

type DriverService interface {
	UpdateDriverLocation(driverID string, request dto.UpdateLocationRequest) (dto.DriverResponse, error)
	SearchDrivers(request dto.SearchDriversRequest) ([]dto.DriverSearchResponse, error)
}

type driverService struct {
	repo repository.DriverRepository
}

func NewDriverService(repo repository.DriverRepository) DriverService {
	return &driverService{repo: repo}
}

func (s *driverService) UpdateDriverLocation(driverID string, request dto.UpdateLocationRequest) (dto.DriverResponse, error) {
	// Parse UUID
	id, err := uuid.Parse(driverID)
	if err != nil {
		return dto.DriverResponse{}, err
	}

	// Update location in repository
	driver, err := s.repo.UpdateLocation(id, request.Latitude, request.Longitude)
	if err != nil {
		return dto.DriverResponse{}, err
	}

	// Convert to DTO
	response := dto.ToDriverResponse(
		driver.ID,
		driver.Name,
		driver.Phone,
		driver.Status,
		driver.Latitude,
		driver.Longitude,
	)

	return response, nil
}

func (s *driverService) SearchDrivers(request dto.SearchDriversRequest) ([]dto.DriverSearchResponse, error) {
	// Set default radius if not provided
	radius := request.Radius
	if radius == 0 {
		radius = 5.0 // Default 5km
	}

	// Search nearby drivers
	drivers, err := s.repo.SearchNearbyDrivers(request.Latitude, request.Longitude, radius)
	if err != nil {
		return nil, err
	}

	return drivers, nil
}
