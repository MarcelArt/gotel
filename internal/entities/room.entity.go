package entities

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	RoomNumber string `gorm:"not null;unique" json:"roomNumber"`
	Floor      string `gorm:"not null" json:"floor"`
	Status     string `gorm:"not null" json:"status"` // VACANT, OCCUPIED, DIRTY, CLEANING, OUT_OF_ORDER
}
