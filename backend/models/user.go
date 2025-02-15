package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name      string               `bson:"name" json:"name"`
	Email     string               `bson:"email" json:"email"`
	Password  string               `bson:"password" json:"-"`
	Followers []primitive.ObjectID `bson:"followers,omitempty" json:"followers,omitempty"`
	Following []primitive.ObjectID `bson:"following,omitempty" json:"following,omitempty"`
	CreatedAt time.Time            `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt,omitempty" json:"updatedAt"`
}
