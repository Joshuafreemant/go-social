package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId      string             `json:"user_id,omitempty" bson:"user_id,omitempty" binding:"required"`
	Content     string             `json:"content,omitempty" bson:"content,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Impressions int                `json:"impressions,omitempty" bson:"impressions,omitempty"`
	Likes       int                `json:"likes,omitempty" bson:"likes,omitempty"`
	Comments    []Comment          `json:"comments,omitempty" bson:"comments,omitempty"`
	Activities  []Activity         `json:"activities,omitempty" bson:"activities,omitempty"`
	Images      []string           `json:"images,omitempty" bson:"images,omitempty"`
}

type Comment struct {
	UserId    string    `json:"user_id,omitempty" bson:"user_id,omitempty" binding:"required"`
	Comment   string    `json:"comment,omitempty" bson:"comment,omitempty" binding:"required"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type Activity struct {
	UserId    string    `json:"user_id,omitempty" bson:"user_id,omitempty" binding:"required"`
	Type      string    `json:"type,omitempty" bson:"type,omitempty"` // e.g., "like", "comment"
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
