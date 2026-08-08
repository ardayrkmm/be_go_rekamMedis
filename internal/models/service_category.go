package models

import (
	"time"
)

type ServiceCategory struct {
	ID        string     `firestore:"id,omitempty" json:"id"`
	Name      string     `json:"name" firestore:"name"`
	CreatedAt time.Time  `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" firestore:"updated_at"`
	DeletedAt *time.Time `json:"-" firestore:"DeletedAt"`
}
