package model

import (
	"time"

	"github.com/google/uuid"
)

type RatingImage struct {
	ID        string    `bson:"_id" json:"id"`
	RatingID  string    `bson:"rating_id" json:"rating_id"`
	Image     string    `bson:"image" json:"image"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (r *RatingImage) BeforeCreate() {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now()
	}
	if r.UpdatedAt.IsZero() {
		r.UpdatedAt = time.Now()
	}
}
