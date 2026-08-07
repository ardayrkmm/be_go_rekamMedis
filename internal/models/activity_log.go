package models

import (
	"time"
)

type ActivityLog struct {
	ID string `firestore:"id,omitempty" json:"id"`
	LogName     string    `json:"log_name"`
	Description string    `json:"description"`
	SubjectType string    `json:"subject_type"`
	SubjectID   *int     `json:"subject_id"`
	CauserType  string    `json:"causer_type"`
	CauserID    *int     `json:"causer_id"`
	Properties  string    `json:"properties"` // MySQL JSON column
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
