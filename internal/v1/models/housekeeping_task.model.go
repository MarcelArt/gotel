package models

import (
	"time"

	"github.com/MarcelArt/gotel/internal/common"
)

type HousekeepingTaskInput struct {
	common.InputModel
	Priority    uint       `gorm:"not null" json:"priority"`
	StartedAt   *time.Time `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt"`
	Note        string     `json:"note"`
	AssigneeID  uint       `gorm:"not null" json:"assigneeId"`
	AssignerID  uint       `gorm:"not null" json:"assignerId"`
	RoomID      uint       `gorm:"not null" json:"roomId"`
}

type HousekeepingTaskPage struct {
	ID          uint       `json:"ID"`
	Priority    uint       `json:"priority"`
	StartedAt   *time.Time `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt"`
	Note        string     `json:"note"`
	AssigneeID  uint       `json:"assigneeId"`
	AssignerID  uint       `json:"assignerId"`
	RoomID      uint       `json:"roomId"`
}
