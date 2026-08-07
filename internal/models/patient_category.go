package models

import (
	"time"
)

type PatientCategory struct {
	ID string `firestore:"id,omitempty" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
