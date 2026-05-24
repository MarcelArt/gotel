package inventory_transactions

import (
	"time"

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
	ID              uint      `json:"ID"`
	CreatedAt       time.Time `json:"createdAt"`
	TransactionType string    `json:"transactionType"`
	Quantity        float64   `json:"quantity"`
	Note            string    `json:"note"`
	ItemID          uint      `json:"itemId"`
	Item            string    `json:"item"`
	Unit            string    `json:"unit"`
	Actor           string    `json:"actor"`
	From            string    `json:"from"`
	To              string    `json:"to"`
}

type ItemCount struct {
	TransactionType string  `json:"transactionType"`
	Quantity        float64 `json:"quantity"`
}
