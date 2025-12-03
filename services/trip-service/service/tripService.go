package service

import (
	"trip-service/dto"
	apperror "trip-service/error"
	"trip-service/model"
	"trip-service/repository"

	"github.com/google/uuid"
)

type TripService interface {
	CreateTrip(request dto.CreateTripRequest) (dto.TripResponse, error)
	GetTrip(tripID string) (dto.TripResponse, error)
	CancelTrip(tripID string, request dto.CancelTripRequest) (dto.TripResponse, error)
	UpdateTripStatus(tripID string, request dto.UpdateTripStatusRequest) (dto.TripResponse, error)
}

type tripService struct {
	repo repository.TripRepository
}

func NewTripService(repo repository.TripRepository) TripService {
	return &tripService{repo: repo}
}

func (s *tripService) CreateTrip(request dto.CreateTripRequest) (dto.TripResponse, error) {
	// Parse passenger ID
	passengerID, err := uuid.Parse(request.PassengerID)
	if err != nil {
		return dto.TripResponse{}, apperror.NewAppError(400, "Invalid passenger ID")
	}

	// Create trip model
	trip := &model.Trip{
		PassengerID:      passengerID,
		PickupLatitude:   request.PickupLatitude,
		PickupLongitude:  request.PickupLongitude,
		DropoffLatitude:  request.DropoffLatitude,
		DropoffLongitude: request.DropoffLongitude,
		Status:           model.StatusSearchingDriver,
	}

	// Save to database
	err = s.repo.CreateTrip(trip)
	if err != nil {
		return dto.TripResponse{}, apperror.NewAppErrorWithErr(500, "Failed to create trip", err)
	}

	// Convert to DTO
	response := dto.ToTripResponse(
		trip.ID,
		trip.PassengerID,
		trip.DriverID,
		trip.PickupLatitude,
		trip.PickupLongitude,
		trip.DropoffLatitude,
		trip.DropoffLongitude,
		trip.Status,
		trip.Fare,
		trip.CancelReason,
		trip.CreatedAt,
		trip.UpdatedAt,
	)

	return response, nil
}

func (s *tripService) GetTrip(tripID string) (dto.TripResponse, error) {
	// Parse trip ID
	id, err := uuid.Parse(tripID)
	if err != nil {
		return dto.TripResponse{}, apperror.NewAppError(400, "Invalid trip ID")
	}

	// Get trip from database
	trip, err := s.repo.GetTripByID(id)
	if err != nil {
		return dto.TripResponse{}, apperror.NewAppErrorWithErr(404, "Trip not found", err)
	}

	// Convert to DTO
	response := dto.ToTripResponse(
		trip.ID,
		trip.PassengerID,
		trip.DriverID,
		trip.PickupLatitude,
		trip.PickupLongitude,
		trip.DropoffLatitude,
		trip.DropoffLongitude,
		trip.Status,
		trip.Fare,
		trip.CancelReason,
		trip.CreatedAt,
		trip.UpdatedAt,
	)

	return response, nil
}

func (s *tripService) CancelTrip(tripID string, request dto.CancelTripRequest) (dto.TripResponse, error) {
	// Parse trip ID
	id, err := uuid.Parse(tripID)
	if err != nil {
		return dto.TripResponse{}, apperror.NewAppError(400, "Invalid trip ID")
	}

	// Get trip from database
	trip, err := s.repo.GetTripByID(id)
	if err != nil {
		return dto.TripResponse{}, apperror.NewAppErrorWithErr(404, "Trip not found", err)
	}

	// Check if trip can be cancelled
	if trip.Status == model.StatusCompleted {
		return dto.TripResponse{}, apperror.NewAppError(400, "Cannot cancel completed trip")
	}

	if trip.Status == model.StatusCancelled {
		return dto.TripResponse{}, apperror.NewAppError(400, "Trip is already cancelled")
	}

	// Update trip status
	trip.Status = model.StatusCancelled
	if request.Reason != "" {
		trip.CancelReason = &request.Reason
	}

	err = s.repo.UpdateTrip(trip)
	if err != nil {
		return dto.TripResponse{}, apperror.NewAppErrorWithErr(500, "Failed to cancel trip", err)
	}

	// Convert to DTO
	response := dto.ToTripResponse(
		trip.ID,
		trip.PassengerID,
		trip.DriverID,
		trip.PickupLatitude,
		trip.PickupLongitude,
		trip.DropoffLatitude,
		trip.DropoffLongitude,
		trip.Status,
		trip.Fare,
		trip.CancelReason,
		trip.CreatedAt,
		trip.UpdatedAt,
	)

	return response, nil
}

func (s *tripService) UpdateTripStatus(tripID string, request dto.UpdateTripStatusRequest) (dto.TripResponse, error) {
	// Parse trip ID
	id, err := uuid.Parse(tripID)
	if err != nil {
		return dto.TripResponse{}, apperror.NewAppError(400, "Invalid trip ID")
	}

	// Get trip from database
	trip, err := s.repo.GetTripByID(id)
	if err != nil {
		return dto.TripResponse{}, apperror.NewAppErrorWithErr(404, "Trip not found", err)
	}

	// Update status
	trip.Status = request.Status

	// Update driver ID if provided
	if request.DriverID != nil {
		driverID, err := uuid.Parse(*request.DriverID)
		if err != nil {
			return dto.TripResponse{}, apperror.NewAppError(400, "Invalid driver ID")
		}
		trip.DriverID = &driverID
	}

	// Update fare if provided
	if request.Fare != nil {
		trip.Fare = request.Fare
	}

	err = s.repo.UpdateTrip(trip)
	if err != nil {
		return dto.TripResponse{}, apperror.NewAppErrorWithErr(500, "Failed to update trip", err)
	}

	// Convert to DTO
	response := dto.ToTripResponse(
		trip.ID,
		trip.PassengerID,
		trip.DriverID,
		trip.PickupLatitude,
		trip.PickupLongitude,
		trip.DropoffLatitude,
		trip.DropoffLongitude,
		trip.Status,
		trip.Fare,
		trip.CancelReason,
		trip.CreatedAt,
		trip.UpdatedAt,
	)

	return response, nil
}
