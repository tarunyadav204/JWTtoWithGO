package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ToDoList struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email,omitempty"`
	Title      string             `json:"title,omitempty"`
	Content    string              `json:"content,omitempty"`
	Completed  *bool               `json:"completed,omitempty"`
}
