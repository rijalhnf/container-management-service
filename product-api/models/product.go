package models

import "time"

type Product struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:varchar(255)"`
	Price     int    `gorm:"type:int"`
	Stock     int    `gorm:"type:int"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
