package categories

import (
	"github.com/MarcelArt/gotel/internal/common"
)

type CategoryInput struct {
	common.InputModel
	Value string `gorm:"not null;unique" json:"value"`
}

type CategoryPage struct {
	ID    uint   `json:"ID"`
	Value string `json:"value"`
}
