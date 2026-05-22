package user_roles

import (
	"github.com/MarcelArt/gotel/internal/common"
)

type UserRoleInput struct {
	common.InputModel
	UserID uint `gorm:"not null" json:"userId"`
	RoleID uint `gorm:"not null" json:"roleId"`
}

type UserRolePage struct {
	ID       uint   `json:"ID"`
	UserID   uint   `json:"userId"`
	RoleID   uint   `json:"roleId"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
