package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	WalletAddress string            `json:"wallet_address" validate:"required"`
	Status    	  *string           `json:"status" validate:"required,eq=active|eq=inactive"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
}