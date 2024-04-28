package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//
const TaskCollection = "task"

//
type Task struct {
	ID			primitive.ObjectID	`bson:"_id,omitempty"`
	Title		string				`bson:"title"`
	Description	string        		`bson:"description"`
}