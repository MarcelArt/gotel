package entities

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Code         string   `gorm:"not null;unique" json:"code"`
	Name         string   `gorm:"not null" json:"name"`
	Picture      string   `json:"picture"`
	TrackingMode string   `gorm:"not null" json:"trackingMode"` // CONSUMABLE, REUSABLE, SERIALIZED
	Unit         string   `gorm:"not null" json:"unit"`
	CategoryID   uint     `gorm:"not null" json:"categoryId"`
	Category     Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
}
