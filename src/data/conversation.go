package data

import "go.mongodb.org/mongo-driver/bson/primitive"

type Conversation struct {
	Title     string    `json:"title" bson:"title"`
	Tags      []string  `json:"tags" bson:"tags"`
	Documents []string  `json:"documents" bson:"documents"`
	Messages  []Message `json:"messages" bson:"messages"`
}

type CreateConversationRequest struct {
	Username  string    `json:"username" bson:"username"`
	Title     string    `json:"title" bson:"title"`
	Tags      []string  `json:"tags" bson:"tags"`
	Documents []string  `json:"documents" bson:"documents"`
	Messages  []Message `json:"messages" bson:"messages"`
}

type LoadConversationRequest struct {
	Username     string              `json:"username"`
	PrevObjectID *primitive.ObjectID `json:"prevId"`
	ObjectID     primitive.ObjectID  `json:"id"`
}

type GetConversationDetailsRequest struct {
	Username     string              `json:"username"`
	ObjectID     primitive.ObjectID  `json:"id"`
}

type DeleteConversationRequest struct {
	Username string             `json:"username"`
	ObjectID primitive.ObjectID `json:"id"`
}
