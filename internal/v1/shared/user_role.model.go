package shared

import "github.com/MarcelArt/gotel/internal/common"

type UserRoleInput struct {
	common.InputModel
	UserID uint `gorm:"not null" json:"userId"`
	RoleID uint `gorm:"not null" json:"roleId"`
}
