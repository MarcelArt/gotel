package entities

import "gorm.io/gorm"

type AssetInstance struct {
	gorm.Model
	Code   string `gorm:"not null;unique" json:"code"`
	ItemID uint   `gorm:"not null" json:"itemId"`
	Item   Item   `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
}
