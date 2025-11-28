package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          string    `bson:"_id" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Status      string    `bson:"status" json:"status"`
	SellerID    string    `bson:"seller_id" json:"seller_id"`
	Rating      float64   `bson:"rating" json:"rating"`
	RateCount   int       `bson:"rate_count" json:"rate_count"`
	SoldCount   int       `bson:"sold_count" json:"sold_count"`
	IsActive    bool      `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`

	// Many-to-many relationships stored as arrays of IDs
	CategoryIDs       []string `bson:"category_ids" json:"category_ids,omitempty"`
	SellerCategoryIDs []string `bson:"seller_category_ids" json:"seller_category_ids,omitempty"`
}

// BeforeCreate generates a new UUID for the ID field if not set
func (p *Product) BeforeCreate() {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now()
	}
}
