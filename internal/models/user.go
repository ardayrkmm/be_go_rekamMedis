package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type RoleEnum string

const (
	RoleAdmin RoleEnum = "admin"
	RoleOwner RoleEnum = "owner"
	RoleStaff RoleEnum = "staff"
	RolePasien RoleEnum = "pasien"
	RoleFisioterapis RoleEnum = "fisioterapis"
)

type User struct {
	ID string `firestore:"id,omitempty" json:"id" gorm:"type:varchar(36);primaryKey"`
	Name      string         `json:"name"`
	Email     string         `json:"email" gorm:"uniqueIndex;type:varchar(255)"`
	Password  string         `json:"-"`
	Role      string         `json:"role"`
	Photo     *string        `json:"photo"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time `json:"-" firestore:"DeletedAt"`
}

func (u *User) GetPhotoUrl() *string {
	if u.Photo != nil {
		url := "storage/" + *u.Photo
		return &url
	}
	return nil
}


func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
