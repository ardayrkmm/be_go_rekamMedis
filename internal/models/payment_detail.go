package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type PaymentDetail struct {
	ID string `firestore:"id,omitempty" json:"id" gorm:"type:varchar(36);primaryKey"`
	PaymentID       string  `json:"payment_id" gorm:"type:varchar(36);index"`
	ServiceMasterID string  `json:"service_master_id,omitempty" firestore:"service_master_id" gorm:"type:varchar(36);index"`
	ItemName    string         `json:"item_name"`
	Quantity    int            `json:"quantity"`
	Price       float64        `json:"price"`
	Subtotal    float64        `json:"subtotal"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt *time.Time `json:"-" firestore:"DeletedAt"`

	Payment *Payment `json:"payment,omitempty" gorm:"-"`
}


func (m *PaymentDetail) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
