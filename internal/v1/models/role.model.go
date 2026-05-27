package models

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/pkg/jsonb"
)

type RoleInput struct {
	common.InputModel
	Name        string                `gorm:"not null;unique" json:"name"`
	Description string                `json:"description"`
	Permissions jsonb.JSONB[[]string] `json:"permissions"`
}

type RolePage struct {
	ID          uint                  `json:"ID"`
	Name        string                `gorm:"not null;unique" json:"name"`
	Description string                `json:"description"`
	Permissions jsonb.JSONB[[]string] `json:"permissions"`
}
