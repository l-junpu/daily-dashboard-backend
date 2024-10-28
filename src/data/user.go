package data

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Username string `json:"username"`
}

type UserDetails struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MongoConvoDetails struct {
	Title    string             `json:"title"`
	ObjectID primitive.ObjectID `json:"id"`
}

type MongoUserDetails struct {
	Username      string              `bson:"username"`
	Password      string              `bson:"password"`
	Conversations []MongoConvoDetails `bson:"conversations"`
}
