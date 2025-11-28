package model

import (
	"time"

	"github.com/google/uuid"
)

type ProductImages struct {
	ID        string    `bson:"_id" json:"id"`
	ProductID string    `bson:"product_id" json:"product_id"`
	Image     string    `bson:"image" json:"image"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (p *ProductImages) BeforeCreate() {
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
