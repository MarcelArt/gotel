package entities

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Value       string `gorm:"not null;unique" json:"value"`
	Description string `json:"description"`
}
