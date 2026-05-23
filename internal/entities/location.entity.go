package entities

import "gorm.io/gorm"

type Location struct {
	gorm.Model
	Value       string `gorm:"not null" json:"value"`
	Description string `json:"description"`
	IsVirtual   bool   `gorm:"not null;default:false" json:"isVirtual"`
}
