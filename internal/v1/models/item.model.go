package models

import (
	"github.com/MarcelArt/gotel/internal/common"
)

type ItemInput struct {
	common.InputModel
	Code         string `gorm:"not null;unique" json:"code"`
	Name         string `gorm:"not null" json:"name"`
	Picture      string `json:"picture"`
	TrackingMode string `gorm:"not null" json:"trackingMode"` // CONSUMABLE, REUSABLE, SERIALIZED
	Unit         string `gorm:"not null" json:"unit"`
	CategoryID   uint   `gorm:"not null" json:"categoryId"`
}

type ItemPage struct {
	ID           uint   `json:"ID"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Picture      string `json:"picture"`
	TrackingMode string `json:"trackingMode"`
	Unit         string `json:"unit"`
	CategoryID   uint   `json:"categoryId"`
}
