package models

import (
	"time"
)

type PaymentDetail struct {
	ID string `firestore:"id,omitempty" json:"id"`
	PaymentID string `json:"payment_id"`
	ServiceMasterID string `json:"service_master_id,omitempty" firestore:"service_master_id"`
	ItemName    string         `json:"item_name"`
	Quantity    int            `json:"quantity"`
	Price       float64        `json:"price"`
	Subtotal    float64        `json:"subtotal"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt *time.Time `json:"-" firestore:"DeletedAt"`

	Payment *Payment `json:"payment,omitempty"`
}
