package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Transaction holds details about a specific transaction.
type Transaction struct {
	TransactionHash      string    `json:"transaction_hash" bson:"transaction_hash" validate:"required"`
	WalletAddress        string    `json:"wallet_address" bson:"wallet_address" validate:"required"`
	TransactionCompleted time.Time `json:"transaction_completed" bson:"transaction_completed"`
	GasFeeUsed           float64   `json:"gas_fee_used" bson:"gas_fee_used"`
}

// OrderItem represents an order with its details and associated transactions.
type OrderItem struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OrderID         string             `json:"order_id" bson:"order_id" validate:"required"`
	OrderPlacedTime time.Time          `json:"order_placed_time" bson:"order_placed_time"`
	OrderCompletedTime time.Time        `json:"order_completed_time" bson:"order_completed_time"`
	OrderStatus     string             `json:"order_status" bson:"order_status" validate:"required,eq=active|eq=inactive|eq=pending|eq=completed"`
	WalletAddress   string             `json:"wallet_address" bson:"wallet_address" validate:"required"`
	GasCollected    float64            `json:"gas_collected" bson:"gas_collected"`
	GasUsed         float64            `json:"gas_used" bson:"gas_used"`
	Addresses       []Transaction      `json:"addresses" bson:"addresses"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}
