package users

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
)

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

type LoginInput struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	IsRemember bool   `json:"isRemember"`
}

type LoginResponse struct {
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
	User         entities.User `json:"user"`
}
