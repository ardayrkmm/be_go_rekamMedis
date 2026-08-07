package models

import (
	"time"
)

type ServiceMaster struct {
	ID          string     `firestore:"id,omitempty" json:"id"`
	Code        string     `json:"code" firestore:"code"`
	Name        string     `json:"name"`
	Category    string     `json:"category" firestore:"category"`
	Duration    int        `json:"duration" firestore:"duration"`
	Description string     `json:"description"`
	BasePrice   float64    `json:"price" firestore:"base_price"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" firestore:"DeletedAt"`
}
