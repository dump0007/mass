package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


// OrderItem represents an order with its details and associated transactions.
type Order struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OrderID         float64            `json:"order_id" bson:"order_id" validate:"required"`
	OrderStatus     string             `json:"order_status" bson:"order_status" validate:"required,eq=active|eq=inactive|eq=pending|stopped|eq=completed"`
	OrderType     	string             `json:"order_type" bson:"order_type" validate:"eq=msg|eq=nft"`
	WalletAddress   string             `json:"wallet_address" bson:"wallet_address" validate:"required"`
	GasCollected    float64            `json:"gas_collected" bson:"gas_collected"`
	GasUsed         float64            `json:"gas_used" bson:"gas_used"`
	Nonce         	float64            `json:"nonce" bson:"nonce"`
	FileName		string			   `json:"file_name" bson:"file_name" validate:"required"`
	Ipfs			string			   `json:"ipfs" bson:"ipfs"`
	PrivateKey		string			   `json:"pvt_key" bson:"pvt_key"`
	Addresses       []Transaction      `json:"addresses" bson:"addresses"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}
