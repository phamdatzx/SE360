package repository

import (
	"driver-service/dto"
	"driver-service/model"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DriverRepository interface {
	UpdateLocation(driverID uuid.UUID, latitude, longitude float64) (*model.Driver, error)
	GetDriverByID(driverID uuid.UUID) (*model.Driver, error)
	SearchNearbyDrivers(latitude, longitude, radius float64) ([]dto.DriverSearchResponse, error)
	CreateDriver(driver *model.Driver) error
	SeedMockData() error
}

type driverRepository struct {
	db *gorm.DB
}

func NewDriverRepository(db *gorm.DB) DriverRepository {
	return &driverRepository{db: db}
}

func (r *driverRepository) UpdateLocation(driverID uuid.UUID, latitude, longitude float64) (*model.Driver, error) {
	var driver model.Driver
	err := r.db.First(&driver, "id = ?", driverID).Error
	if err != nil {
		return nil, err
	}

	driver.Latitude = latitude
	driver.Longitude = longitude

	err = r.db.Save(&driver).Error
	if err != nil {
		return nil, err
	}

	return &driver, nil
}

func (r *driverRepository) GetDriverByID(driverID uuid.UUID) (*model.Driver, error) {
	var driver model.Driver
	err := r.db.First(&driver, "id = ?", driverID).Error
	if err != nil {
		return nil, err
	}
	return &driver, nil
}

func (r *driverRepository) SearchNearbyDrivers(latitude, longitude, radius float64) ([]dto.DriverSearchResponse, error) {
	var results []dto.DriverSearchResponse

	// Haversine formula to calculate distance
	// 6371 is Earth's radius in kilometers
	query := `
		SELECT 
			id,
			name,
			phone,
			status,
			latitude,
			longitude,
			(6371 * acos(
				cos(radians(?)) * cos(radians(latitude)) * 
				cos(radians(longitude) - radians(?)) + 
				sin(radians(?)) * sin(radians(latitude))
			)) AS distance
		FROM drivers
		WHERE status = ?
		HAVING distance < ?
		ORDER BY distance
		LIMIT 50
	`

	err := r.db.Raw(query, latitude, longitude, latitude, model.StatusAvailable, radius).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *driverRepository) CreateDriver(driver *model.Driver) error {
	return r.db.Create(driver).Error
}

func (r *driverRepository) SeedMockData() error {
	// Check if data already exists
	var count int64
	r.db.Model(&model.Driver{}).Count(&count)
	if count > 0 {
		fmt.Println("Mock data already exists, skipping seed")
		return nil
	}

	// Central point: Ho Chi Minh City center (approximately)
	centerLat := 10.762622
	centerLng := 106.660172

	mockDrivers := []model.Driver{
		{Name: "Nguyễn Văn A", Phone: "+84901234567", Status: model.StatusAvailable, Latitude: centerLat + 0.01, Longitude: centerLng + 0.01},
		{Name: "Trần Văn B", Phone: "+84901234568", Status: model.StatusAvailable, Latitude: centerLat - 0.01, Longitude: centerLng - 0.01},
		{Name: "Lê Văn C", Phone: "+84901234569", Status: model.StatusAvailable, Latitude: centerLat + 0.02, Longitude: centerLng - 0.02},
		{Name: "Phạm Văn D", Phone: "+84901234570", Status: model.StatusAvailable, Latitude: centerLat - 0.02, Longitude: centerLng + 0.02},
		{Name: "Hoàng Văn E", Phone: "+84901234571", Status: model.StatusAvailable, Latitude: centerLat + 0.015, Longitude: centerLng + 0.015},
		{Name: "Võ Văn F", Phone: "+84901234572", Status: model.StatusBusy, Latitude: centerLat - 0.015, Longitude: centerLng - 0.015},
		{Name: "Đặng Văn G", Phone: "+84901234573", Status: model.StatusAvailable, Latitude: centerLat + 0.03, Longitude: centerLng + 0.03},
		{Name: "Bùi Văn H", Phone: "+84901234574", Status: model.StatusAvailable, Latitude: centerLat - 0.03, Longitude: centerLng - 0.03},
		{Name: "Đỗ Văn I", Phone: "+84901234575", Status: model.StatusOffline, Latitude: centerLat + 0.005, Longitude: centerLng + 0.005},
		{Name: "Ngô Văn K", Phone: "+84901234576", Status: model.StatusAvailable, Latitude: centerLat - 0.005, Longitude: centerLng - 0.005},
		{Name: "Dương Văn L", Phone: "+84901234577", Status: model.StatusAvailable, Latitude: centerLat + 0.025, Longitude: centerLng - 0.025},
		{Name: "Hồ Văn M", Phone: "+84901234578", Status: model.StatusAvailable, Latitude: centerLat - 0.025, Longitude: centerLng + 0.025},
		{Name: "Vũ Văn N", Phone: "+84901234579", Status: model.StatusAvailable, Latitude: centerLat + 0.008, Longitude: centerLng - 0.008},
		{Name: "Phan Văn O", Phone: "+84901234580", Status: model.StatusBusy, Latitude: centerLat - 0.008, Longitude: centerLng + 0.008},
		{Name: "Lý Văn P", Phone: "+84901234581", Status: model.StatusAvailable, Latitude: centerLat + 0.012, Longitude: centerLng + 0.012},
	}

	// Add some randomness to locations
	for i := range mockDrivers {
		mockDrivers[i].Latitude += (rand.Float64() - 0.5) * 0.001
		mockDrivers[i].Longitude += (rand.Float64() - 0.5) * 0.001
	}

	// Bulk insert
	result := r.db.Create(&mockDrivers)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("✅ Mock data seeded successfully: %d drivers created\n", len(mockDrivers))
	return nil
}
