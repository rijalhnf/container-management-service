package models

import "time"

type Voyage struct {
	ID uint `gorm:"primaryKey" json:"id"`

	VoyageNumber string `gorm:"size:50;uniqueIndex;not null" json:"voyage_number"`

	VesselName string `gorm:"size:150;not null" json:"vessel_name"`

	OriginPortID uint `json:"origin_port_id"`
	OriginPort   Port `gorm:"foreignKey:OriginPortID"`

	DestinationPortID uint `json:"destination_port_id"`
	DestinationPort   Port `gorm:"foreignKey:DestinationPortID"`

	ETA time.Time `json:"eta"`

	ETD time.Time `json:"etd"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
