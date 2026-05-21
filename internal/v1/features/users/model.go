package users

import "github.com/MarcelArt/gotel/internal/common"

type UserInput struct {
	common.InputModel
	Username string `gorm:"not null;unique" json:"username"`
	Email    string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null" json:"password"`
}

type UserPage struct {
	ID       uint   `json:"ID"`
	Username string `gorm:"not null;unique" json:"username"`
	Email    string `gorm:"not null;unique" json:"email"`
}
