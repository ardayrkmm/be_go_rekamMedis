package models

import (
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
	ID string `firestore:"id,omitempty" json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
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
