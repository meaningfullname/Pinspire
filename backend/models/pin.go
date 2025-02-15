package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	ID  string `bson:"id" json:"id"`
	URL string `bson:"url" json:"url"`
}

type Comment struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User    primitive.ObjectID `bson:"user" json:"user"`
	Name    string             `bson:"name" json:"name"`
	Comment string             `bson:"comment" json:"comment"`
}

type Pin struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	PinText   string             `bson:"pin" json:"pin"`
	Owner     primitive.ObjectID `bson:"owner" json:"owner"`
	Image     Image              `bson:"image" json:"image"`
	Comments  []Comment          `bson:"comments,omitempty" json:"comments,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}
