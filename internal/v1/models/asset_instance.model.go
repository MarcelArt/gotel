package models

import (
	"github.com/MarcelArt/gotel/internal/common"
)

type AssetInstanceInput struct {
	common.InputModel
	Code   string `gorm:"not null;unique" json:"code"`
	ItemID uint   `gorm:"not null" json:"itemId"`
}

type AssetInstancePage struct {
	ID       uint   `json:"ID"`
	Code     string `json:"code"`
	ItemID   uint   `json:"itemId"`
	Status   string `json:"status"`
	Location string `json:"location"`
	Note     string `json:"note"`
}
