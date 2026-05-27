package models

import (
	"time"

	"github.com/MarcelArt/gotel/internal/common"
)

type AssetTransactionInput struct {
	common.InputModel
	TransactionType string `gorm:"not null" json:"transactionType"` // ASSIGN, REPAIR_IN, REPAIR_OUT, DISPOSE, LOST, BROKEN, TRANSFER, OTHER
	Status          string `gorm:"not null" json:"status"`          // AVAILABLE, IN_USE, REPAIRING, LOST, DISPOSED, BROKEN
	Note            string `json:"note"`
	LocationID      uint   `gorm:"not null" json:"locationId"`
	InstanceID      uint   `gorm:"not null" json:"instanceId"`
	ActorID         uint   `gorm:"not null" json:"actorId"`
}

type AssetTransactionPage struct {
	ID              uint      `json:"ID"`
	CreatedAt       time.Time `json:"createdAt"`
	TransactionType string    `json:"transactionType"`
	Status          string    `json:"status"`
	Note            string    `json:"note"`
	LocationID      uint      `json:"locationId"`
	InstanceID      uint      `json:"instanceId"`
	ActorID         uint      `json:"actorId"`
	Location        string    `json:"location"`
	Actor           string    `json:"actor"`
}
