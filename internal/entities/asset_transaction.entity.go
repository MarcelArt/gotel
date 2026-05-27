package entities

import "gorm.io/gorm"

type AssetTransaction struct {
	gorm.Model
	TransactionType string        `gorm:"not null" json:"transactionType"` // ASSIGN, REPAIR_IN, REPAIR_OUT, DISPOSE, LOST, BROKEN, TRANSFER, OTHER
	Status          string        `gorm:"not null" json:"status"`          // AVAILABLE, IN_USE, REPAIRING, LOST, DISPOSED, BROKEN
	Note            string        `json:"note"`
	LocationID      uint          `gorm:"not null" json:"locationId"`
	Location        Location      `gorm:"foreignKey:LocationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	InstanceID      uint          `gorm:"not null" json:"instanceId"`
	Instance        AssetInstance `gorm:"foreignKey:InstanceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	ActorID         uint          `gorm:"not null" json:"actorId"`
	Actor           User          `gorm:"foreignKey:ActorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
}
