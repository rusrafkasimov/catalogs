package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Catalog struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Active    bool               `bson:"active" json:"active"`
	Category  string             `bson:"category" json:"category"`
	Name      string             `bson:"name" json:"name"`
	Desc      string             `bson:"desc" json:"desc"`
	Value     string             `bson:"value" json:"value"`

	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
