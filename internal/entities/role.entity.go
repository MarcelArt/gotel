package entities

import (
	"github.com/MarcelArt/gotel/pkg/jsonb"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string                `gorm:"not null;unique" json:"name"`
	Description string                `json:"description"`
	Permissions jsonb.JSONB[[]string] `json:"permissions"`
}
