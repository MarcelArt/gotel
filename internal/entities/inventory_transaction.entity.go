package entities

import "gorm.io/gorm"

type InventoryTransaction struct {
	gorm.Model
	TransactionType string    `gorm:"not null" json:"transactionType"`
	Quantity        float64   `gorm:"not null" json:"quantity"`
	Note            string    `json:"note"`
	ItemID          uint      `gorm:"not null" json:"itemId"`
	FromID          *uint     `json:"fromId"`
	ToID            *uint     `json:"toId"`
	ActorID         uint      `gorm:"not null" json:"actorId"`
	Item            Item      `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	From            *Location `gorm:"foreignKey:FromID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	To              *Location `gorm:"foreignKey:ToID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Actor           User      `gorm:"foreignKey:ActorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
}
