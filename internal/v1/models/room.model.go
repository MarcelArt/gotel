package models

import (
	"time"

	"github.com/MarcelArt/gotel/internal/common"
)

type RoomInput struct {
	common.InputModel
	RoomNumber string `gorm:"not null;unique" json:"roomNumber"`
	Floor      string `gorm:"not null" json:"floor"`
	Status     string `gorm:"not null" json:"status"` // VACANT, OCCUPIED, DIRTY, CLEANING, OUT_OF_ORDER
}

type RoomPage struct {
	ID            uint       `json:"ID"`
	RoomNumber    string     `json:"roomNumber"`
	Floor         string     `json:"floor"`
	Status        string     `json:"status"`
	TaskID        uint       `json:"taskId"`
	TaskStartedAt *time.Time `json:"taskStartedAt"`
	AssigneeID    uint       `json:"assigneeId"`
	Assignee      string     `json:"assignee"`
}
