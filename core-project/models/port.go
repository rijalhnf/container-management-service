package models

import "time"

type Port struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Code string `gorm:"size:10;uniqueIndex;not null" json:"code"`

	Name string `gorm:"size:150;not null" json:"name"`

	Country string `gorm:"size:100" json:"country"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
