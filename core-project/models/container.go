package models

import "time"

type Container struct {
	ID uint `gorm:"primaryKey" json:"id"`

	ContainerNumber string `gorm:"type:varchar(11);uniqueIndex;not null" json:"container_number"`

	ContainerTypeID uint          `json:"container_type_id"`
	ContainerType   ContainerType `gorm:"foreignKey:ContainerTypeID"`

	ShippingLineID uint         `json:"shipping_line_id"`
	ShippingLine   ShippingLine `gorm:"foreignKey:ShippingLineID"`

	PortID uint `json:"port_id"`
	Port   Port `gorm:"foreignKey:PortID"`

	VoyageID uint   `json:"voyage_id"`
	Voyage   Voyage `gorm:"foreignKey:VoyageID"`

	Status string `gorm:"type:varchar(20);default:'ARRIVED'" json:"status"`

	CreatedBy uint `json:"created_by"`
	User      User `gorm:"foreignKey:CreatedBy"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
