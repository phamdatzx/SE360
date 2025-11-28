package model

import (
	"time"

	"github.com/google/uuid"
)

type Rating struct {
	ID          string    `bson:"_id" json:"id"`
	ParentID    *string   `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Score       float64   `bson:"score" json:"score"`
	UserID      string    `bson:"user_id" json:"user_id"`
	Description string    `bson:"description" json:"description"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

func (r *Rating) BeforeCreate() {
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
