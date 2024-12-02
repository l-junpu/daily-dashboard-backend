package llm

import (
	"context"
	"daily-dashboard-backend/src/data"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Note: Conversations need tags

type MongoDBClient struct {
	MongoClient *mongo.Client
}

// Create a MongoDB Client that timeout after 10s of failing to perform a task
func CreateMongoDBClient() (*MongoDBClient, error) {
	mongoUri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUri).SetServerSelectionTimeout(10*time.Second))
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
func (s *MongoDBClient) findUser(username string) (*data.MongoUserDetails, error) {
	filter := bson.M{"username": username}
	coll := s.MongoClient.Database("UserData").Collection("Users")

	var user data.MongoUserDetails
	err := coll.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

// Inserts a new User into MongoDB upon successful registration
func (s *MongoDBClient) InsertUser(username string, password string) error {
	coll := s.MongoClient.Database("UserData").Collection("Users")

	data := data.MongoUserDetails{
		Username:      username,
		Password:      password,
		Conversations: make([]data.MongoConvoDetails, 0),
	}

	bsonUserData, err := bson.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to Marshal MongoUserDetails when inserting user: %w", err)
	}

	user, _ := s.findUser(username)
	if user == nil {
		_, err := coll.InsertOne(context.Background(), bsonUserData)
		if err != nil {
			return fmt.Errorf("unable to Insert MongoUserDetails: %w", err)
		}
	}

	return nil
}

// Updates a User's Conversations[]
func (s *MongoDBClient) UpdateUser(username string, conversations []data.MongoConvoDetails) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"conversations": conversations}}
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

// Add a Conversation Detail to UserData in MongoDB
func (s *MongoDBClient) insertTitle(username string, title string, id interface{}) error {
	filter := bson.M{"username": username}
	userCollection := s.MongoClient.Database("UserData").Collection("Users")

	var user data.MongoUserDetails
	err := userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return fmt.Errorf("unable to find user when inserting title: %w", err)
	}

	// Convert result of InsertOne.InsertedId -> primitive.ObjectId
	objectId, ok := id.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("insertedID is not of type primitive.ObjectID")
	}

	user.Conversations = append(user.Conversations, data.MongoConvoDetails{})
	copy(user.Conversations[1:], user.Conversations)
	user.Conversations[0] = data.MongoConvoDetails{Title: title, ObjectID: objectId}
	return s.UpdateUser(username, user.Conversations)
}

// Removes a Conversation Detail from UserData in MongoDB
func (s *MongoDBClient) deleteTitle(username string, idToRemove primitive.ObjectID) error {
	filter := bson.M{"username": username}
	userCollection := s.MongoClient.Database("UserData").Collection("Users")

	var user data.MongoUserDetails
	err := userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return fmt.Errorf("failed to find user when deleting title: %w", err)
	}

	index := -1
	for i, convoDetails := range user.Conversations {
		if convoDetails.ObjectID == idToRemove {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("attempting to remove a conversationId that is not in UserData")
	}

	// Remove the Conversation Detail by shifting subsequent elements to the index of the idToRemove
	// Update size
	copy(user.Conversations[index:], user.Conversations[index+1:])
	user.Conversations = user.Conversations[:len(user.Conversations)-1]
	s.UpdateUser(user.Username, user.Conversations)

	return nil
}

// Insert Message into Conversation & Bump Title
func (s *MongoDBClient) InsertNewMessage(username string, id primitive.ObjectID, message data.Message) error {
	if err := s.updateConversation(username, id, message); err != nil {
		return err
	}

	err := s.bumpConversationSequence(username, id)
	return err
}

// Retrieves all Conversation Details attached to a User
func (s *MongoDBClient) GetConversationDetails(username string) ([]data.MongoConvoDetails, error) {
	user, err := s.findUser(username)
	if err != nil {
		return make([]data.MongoConvoDetails, 0), fmt.Errorf("unable to find user when retrieving conversation details: %w", err)
	}
	return user.Conversations, nil
}

func (s *MongoDBClient) UpdateConversationDetails(username string, details []data.MongoConvoDetails) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"conversations": details}}
	coll := s.MongoClient.Database("UserData").Collection("Users")

	result, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no documents matched the filter")
	}

	return nil
}

// Bumps the updated conversation to the front of the list
func (s *MongoDBClient) bumpConversationSequence(username string, id primitive.ObjectID) error {
	// Retrieve all Conversations
	convoDetails, err := s.GetConversationDetails(username)
	if err != nil {
		return fmt.Errorf("unable to retrieve conversation details from MongoDB: %w", err)
	}

	// Find the Conversation to be shifted
	var idx int = -1
	for i, detail := range convoDetails {
		if detail.ObjectID == id {
			idx = i
			break
		}
	}

	// If it exists (Which it should...), shift it to the front
	if idx != -1 {
		var updatedConvoDetails = append([]data.MongoConvoDetails{}, convoDetails[:idx]...)
		updatedConvoDetails = append(updatedConvoDetails, convoDetails[idx+1:]...)
		updatedConvoDetails = append([]data.MongoConvoDetails{convoDetails[idx]}, updatedConvoDetails...)

		return s.UpdateConversationDetails(username, updatedConvoDetails)
	} else {
		return fmt.Errorf("unable to find the conversation that was tied to ObjectId")
	}
}

// Retrieves the Actual Conversation from MongoDB
func (s *MongoDBClient) FindConversation(username string, id primitive.ObjectID) (*data.Conversation, error) {
	filter := bson.M{"_id": id}
	coll := s.MongoClient.Database("ConversationData").Collection(username)

	var convo data.Conversation
	err := coll.FindOne(context.Background(), filter).Decode(&convo)
	if err != nil {
		return nil, fmt.Errorf("failed to find conversation: %w", err)
	}

	return &convo, nil
}

// Updates a Conversation's history in MongoDB
func (s *MongoDBClient) updateConversation(username string, id primitive.ObjectID, message data.Message) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$push": bson.M{"messages": message}}
	coll := s.MongoClient.Database("ConversationData").Collection(username)

	result, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no documents matched the filter")
	}

	return nil
}

// Insert a New Conversation in MongoDB
func (s *MongoDBClient) InsertNewConversation(username string, conversation data.Conversation) (primitive.ObjectID, error) {
	bsonConvoData, err := bson.Marshal(conversation)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("unable to Marshal Conversation when inserting convo: %w", err)
	}

	user, _ := s.findUser(username)
	if user != nil {
		coll := s.MongoClient.Database("ConversationData").Collection(username)
		result, err := coll.InsertOne(context.Background(), bsonConvoData)
		if err != nil {
			return primitive.NilObjectID, fmt.Errorf("unable to Insert New Conversation to ConversationData: %w", err)
		}
		err = s.insertTitle(username, conversation.Title, result.InsertedID)
		if err != nil {
			return primitive.NilObjectID, fmt.Errorf("unable to Append New Title to user to UserData: %w", err)
		}

		return result.InsertedID.(primitive.ObjectID), nil
	}

	return primitive.NilObjectID, fmt.Errorf("unable to find User")
}

// Delete a Conversation from MongoDB
func (s *MongoDBClient) DeleteConversation(username string, objectId primitive.ObjectID) error {
	// // This portion removes the Conversation from the DB - Ideally we avoid
	// // this so we can still have our "accurate" summary logs
	// filter := bson.M{"_id": objectId}
	// convoCollection := s.MongoClient.Database("ConversationData").Collection(username)

	// // Removes the Conversation from Username
	// result, err := convoCollection.DeleteOne(context.Background(), filter)
	// if err != nil {
	// 	return fmt.Errorf("failed to delete conversation from ConversationData: %w", err)
	// } else if result.DeletedCount != 1 {
	// 	return fmt.Errorf("expected 1 document to be deleted, got %d", result.DeletedCount)
	// }

	// Removes the TitleId attached to the UserData
	err := s.deleteTitle(username, objectId)
	if err != nil {
		return fmt.Errorf("unable to delete title from UserData %w", err)
	}

	return nil
}
