package models

import (
	"github.com/MarcelArt/gotel/internal/common"
)

type CategoryInput struct {
	common.InputModel
	Value       string `gorm:"not null;unique" json:"value"`
	Description string `json:"description"`
}

type CategoryPage struct {
	ID          uint   `json:"ID"`
	Value       string `json:"value"`
	Description string `json:"description"`
}
