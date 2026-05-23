package inventory_transactions

import (
	"github.com/MarcelArt/gotel/internal/common"
)

type InventoryTransactionInput struct {
	common.InputModel
	TransactionType string  `gorm:"not null" json:"transactionType"`
	Quantity        float64 `gorm:"not null" json:"quantity"`
	Note            string  `json:"note"`
	ItemID          uint    `gorm:"not null" json:"itemId"`
	FromID          *uint   `json:"fromId"`
	ToID            *uint   `json:"toId"`
	ActorID         uint    `gorm:"not null" json:"actorId"`
}

type InventoryTransactionPage struct {
	ID              uint    `json:"ID"`
	TransactionType string  `json:"transactionType"`
	Quantity        float64 `json:"quantity"`
	Note            string  `json:"note"`
	ItemID          uint    `json:"itemId"`
	FromID          *uint   `json:"fromId"`
	ToID            *uint   `json:"toId"`
	ActorID         uint    `json:"actorId"`
}
