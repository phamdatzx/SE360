package model

import (
	"time"

	"github.com/google/uuid"
)

type ProductReport struct {
	ID        string    `bson:"_id" json:"id"`
	ProductID string    `bson:"product_id" json:"product_id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	Type      string    `bson:"type" json:"type"`
	Content   string    `bson:"content" json:"content"`
	Status    string    `bson:"status" json:"status"`
	IsDeleted bool      `bson:"is_deleted" json:"is_deleted"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (p *ProductReport) BeforeCreate() {
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
