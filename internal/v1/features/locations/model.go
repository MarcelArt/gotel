package locations

import (
	"github.com/MarcelArt/gotel/internal/common"
)

type LocationInput struct {
	common.InputModel
	Value       string `gorm:"not null" json:"value"`
	Description string `json:"description"`
	IsVirtual   bool   `gorm:"not null;default:false" json:"isVirtual"`
}

type LocationPage struct {
	ID          uint   `json:"ID"`
	Value       string `json:"value"`
	Description string `json:"description"`
	IsVirtual   bool   `json:"isVirtual"`
}
