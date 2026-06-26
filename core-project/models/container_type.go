package models

import "time"

type ContainerType struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Code string `gorm:"size:10;uniqueIndex;not null" json:"code"`

	Name string `gorm:"size:100;not null" json:"name"`

	Description string `gorm:"size:255" json:"description"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
