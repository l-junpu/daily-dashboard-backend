package llm

import (
	"context"
	"daily-dashboard-backend/src/data"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Note: Conversations need tags

type MongoDBClient struct {
	MongoClient *mongo.Client
}

// Create a MongoDB Client that timeout after 10s of failing to perform a task
func CreateMongoDBClient(uri string) (*MongoDBClient, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri).SetServerSelectionTimeout(10*time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return &MongoDBClient{
		MongoClient: client,
	}, nil
}

// Terminates the MongoDB Connection
func (s *MongoDBClient) Terminate() error {
	return s.MongoClient.Disconnect(context.Background())
}

// Retrieves UserDetails
func (s *MongoDBClient) FindUser(username string) (*data.MongoUserDetails, error) {
	filter := bson.M{"username": username}
	coll := s.MongoClient.Database("UserData").Collection("Users")

	var user data.MongoUserDetails
	err := coll.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

// Updates a User's Titles[]
func (s *MongoDBClient) UpdateUser(username string, titles []string) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"titles": titles}}
	coll := s.MongoClient.Database("UserData").Collection("Users")

	result, err := coll.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.ModifiedCount != 1 {
		return fmt.Errorf("expected 1 document to be modified, got %d", result.ModifiedCount)
	}

	return nil
}

// Retrieves a Conversation from MongoDB
func (s *MongoDBClient) FindConversation(username string, title string) (*data.Conversation, error) {
	filter := bson.M{"title": title}
	coll := s.MongoClient.Database("ConversationData").Collection(username)

	var convo data.Conversation
	err := coll.FindOne(context.Background(), filter).Decode(&convo)
	if err != nil {
		return nil, fmt.Errorf("failed to find conversation: %w", err)
	}

	return &convo, nil
}

// Updates a Conversation in MongoDB
func (s *MongoDBClient) UpdateConversation(username string, convo data.Conversation) error {
	filter := bson.M{"title": convo.Title}
	update := bson.M{"$set": bson.M{"messages": convo.Messages}}
	coll := s.MongoClient.Database("ConversationData").Collection(username)

	result, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	if result.ModifiedCount != 1 {
		return fmt.Errorf("expected 1 document to be modified, got %d", result.ModifiedCount)
	}

	return nil
}

// Inserts a new User into MongoDB upon successful registration
func (s *MongoDBClient) InsertUser(username string, password string) error {
	coll := s.MongoClient.Database("UserData").Collection("Users")

	data := data.MongoUserDetails{
		Username: username,
		Password: password,
		Titles:   make([]string, 0),
	}

	bsonUserData, err := bson.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to Marshal MongoUserDetails when inserting user: %w", err)
	}

	user, _ := s.FindUser(username)
	if user == nil {
		_, err = coll.InsertOne(context.Background(), bsonUserData)
		if err != nil {
			return fmt.Errorf("unable to Insert MongoUserDetails: %w", err)
		}
	}

	return nil
}

// Insert a New Conversation in MongoDB
func (s *MongoDBClient) InsertConversation(username string, conversation data.Conversation) error {
	bsonConvoData, err := bson.Marshal(conversation)
	if err != nil {
		return fmt.Errorf("unable to Marshal Conversation when inserting convo: %w", err)
	}

	user, _ := s.FindUser(username)
	if user != nil {
		coll := s.MongoClient.Database("ConversationData").Collection(username)
		_, err = coll.InsertOne(context.Background(), bsonConvoData)
		if err != nil {
			return fmt.Errorf("unable to Insert New Conversation to ConversationData: %w", err)
		}
		err = s.InsertTitle(username, conversation.Title)
		if err != nil {
			return fmt.Errorf("unable to Append New Title to user to UserData: %w", err)
		}
	}

	return nil
}

// Delete a Conversation from MongoDB
func (s *MongoDBClient) DeleteConversation(username string, title string) error {
	filter := bson.M{"title": title}
	convoCollection := s.MongoClient.Database("ConversationData").Collection(username)

	// Removes the Conversation from Username
	result, err := convoCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete conversation from ConversationData: %w", err)
	}

	// Removes the TitleId attached to the UserData
	err = s.DeleteTitle(username, title)
	if err != nil {
		return fmt.Errorf("unable to delete title from UserData %w", err)
	}

	if result.DeletedCount != 1 {
		return fmt.Errorf("expected 1 document to be deleted, got %d", result.DeletedCount)
	}

	return nil
}

// Retrieves all Titles attached to a User
func (s *MongoDBClient) GetConversationTitles(username string) ([]string, error) {
	user, err := s.FindUser(username)
	if err != nil {
		return make([]string, 0), fmt.Errorf("unable to find user when retrieving titles: %w", err)
	}
	return user.Titles, nil
}

// Adds a Title to UserData in MongoDB
func (s *MongoDBClient) InsertTitle(username string, title string) error {
	filter := bson.M{"username": username}
	userCollection := s.MongoClient.Database("UserData").Collection("Users")

	var user data.MongoUserDetails
	err := userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return fmt.Errorf("unable to find user when inserting title: %w", err)
	}

	user.Titles = append(user.Titles, "")
	copy(user.Titles[1:], user.Titles)
	user.Titles[0] = title
	return s.UpdateUser(username, user.Titles)
}

// Removes a Title from UserData in MongoDB
func (s *MongoDBClient) DeleteTitle(username string, titleToRemove string) error {
	filter := bson.M{"username": username}
	userCollection := s.MongoClient.Database("UserData").Collection("Users")

	var user data.MongoUserDetails
	err := userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return fmt.Errorf("failed to find user when deleting title: %w", err)
	}

	// new
	index := -1
	for i, title := range user.Titles {
		if title == titleToRemove {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("attempting to remove a title that is not in UserData")
	}

	// Remove the Title by shifting subsequent elements to the index of the titleToRemove
	// Update size
	copy(user.Titles[index:], user.Titles[index+1:])
	user.Titles = user.Titles[:len(user.Titles)-1]
	s.UpdateUser(user.Username, user.Titles)

	return nil
}
