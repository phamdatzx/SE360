package model

import (
	"time"

	"github.com/google/uuid"
)

type ProductOption struct {
	ID                string    `bson:"_id" json:"id"`
	ProductID         string    `bson:"product_id" json:"product_id"`
	Name              string    `bson:"name" json:"name"`
	Price             float64   `bson:"price" json:"price"`
	QuantityAvailable int       `bson:"quantity_available" json:"quantity_available"`
	SoldCount         int       `bson:"sold_count" json:"sold_count"`
	Image             string    `bson:"image" json:"image"`
	CreatedAt         time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at" json:"updated_at"`
}

func (p *ProductOption) BeforeCreate() {
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
