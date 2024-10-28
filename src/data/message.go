package data

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	Role    string `json:"role" bson:"role"`
	Content string `json:"content" bson:"content"`
}

type MessageRequest struct {
	Username string             `json:"username"`
	Id       primitive.ObjectID `json:"id"`
	Message  Message            `json:"message"`
}
