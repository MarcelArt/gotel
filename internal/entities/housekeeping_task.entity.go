package entities

import (
	"time"

	"gorm.io/gorm"
)

type HousekeepingTask struct {
	gorm.Model
	Priority    uint       `gorm:"not null" json:"priority"` // Lower number is most urgent
	StartedAt   *time.Time `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt"`
	Note        string     `json:"note"`
	AssigneeID  uint       `gorm:"not null" json:"assigneeId"`
	Assignee    User       `gorm:"foreignKey:AssigneeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	AssignerID  uint       `gorm:"not null" json:"assignerId"`
	Assigner    User       `gorm:"foreignKey:AssignerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	RoomID      uint       `gorm:"not null" json:"roomId"`
	Room        Room       `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
}
