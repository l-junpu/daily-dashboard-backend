package api

import (
	"daily-dashboard-backend/src/data"
	"daily-dashboard-backend/src/inferer"
	"daily-dashboard-backend/src/llm"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HandleGetConvosFromUser(c *llm.MongoDBClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}

		// Decode user data
		var user data.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("unable to decode user details in HandleGetConvosFromUser(): %w", err)
			log.Fatal(err)
		}

		// Retrieve User's Conversation Details
		convoDetails, err := c.GetConversationDetails(user.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("unable to retrieve conversation details from MongoDB: %w", err)
			log.Fatal(err)
		}

		WriteAsJson(w, map[string]interface{}{"convoDetails": convoDetails}, http.StatusOK)
	}
}

func HandleGetConvoDetails(c *llm.MongoDBClient, rc *llm.RedisClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}

		// Decode user data
		var request data.GetConversationDetailsRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("unable to decode user request in HandleGetConvoDetails(): %w", err)
			log.Fatal(err)
		}

		convo, err := c.FindConversation(request.Username, request.ObjectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("unable to retrieve conversation history from MongoDB: %w", err)
			log.Fatal(err)
		}

		WriteAsJson(w, map[string]interface{}{"tags": convo.Tags, "docs": convo.Documents}, http.StatusOK)
	}
}

func HandleGetConvoHistory(c *llm.MongoDBClient, rc *llm.RedisClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}

		// Decode user data
		var request data.LoadConversationRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("unable to decode user request in HandleGetConvoHistory(): %w", err)
			log.Fatal(err)
		}

		convo, err := c.FindConversation(request.Username, request.ObjectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("unable to retrieve conversation history from MongoDB: %w", err)
			log.Fatal(err)
		}

		// Unload the Previous Conversation from Redis
		if request.PrevObjectID != nil {
			if err := rc.RemoveConversationData(request.PrevObjectID.Hex()); err != nil {
				log.Fatal(err)
			}
		}

		// Load the new Conversation into Redis
		if err := rc.SetConversationData(request.ObjectID.Hex(), convo); err != nil {
			log.Fatal(err)
		}

		WriteAsJson(w, map[string]interface{}{"messages": convo.Messages}, http.StatusOK)
	}
}

func HandleCreateNewConvo(c *llm.MongoDBClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}

		// Decode user data
		var request data.CreateConversationRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("unable to decode create convo details in HandleCreateNewConvo(): %w", err)
			log.Fatal(err)
		}

		// Insert a New Conversation into MongoDB
		id, err := c.InsertNewConversation(request.Username, data.Conversation{
			Title:     request.Title,
			Tags:      request.Tags,
			Documents: request.Documents,
			Messages:  request.Messages,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("unable to insert a new conversation into MongoDB: %w", err)
			log.Fatal(err)
		}

		WriteAsJson(w, map[string]interface{}{"id": id.Hex()}, http.StatusOK)
	}
}

func HandleDeleteConvo(c *llm.MongoDBClient, rc *llm.RedisClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}

		// Decode user data
		var request data.DeleteConversationRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("unable to decode delete convo details in HandleDeleteConvo(): %w", err)
			log.Fatal(err)
		}

		// Request MongoDB to "delete" the conversation - In actual fact, we only remove it for the user
		// The actual conversation data still exists for future reference
		if err := c.DeleteConversation(request.Username, request.ObjectID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("unable to delete conversation from MongoDB: %w", err)
			log.Fatal(err)
		}

		// Remove the Conversation from our Redis Client
		if err := rc.RemoveConversationData(request.ObjectID.Hex()); err != nil {
			log.Println(err)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func HandleNewUserPrompt(c *llm.MongoDBClient, rc *llm.RedisClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}

		// Decode user data
		var request data.MessageRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("unable to decode new user prompt details in HandleNewUserPrompt(): %w", err)
			log.Fatal(err)
		}

		// Insert our new prompt
		if err := c.InsertNewMessage(request.Username, request.Id, request.Message); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("unable to add message to conversation in MongoDB: %w", err)
			log.Fatal(err)
		}

		// Grab the conversation from Redis and update it with the New Message
		convo, err := rc.GetConversationData(request.Id.Hex())
		if err != nil {
			log.Fatalln(err)
		}
		convo.Messages = append(convo.Messages, request.Message)

		c.GetConversationDetails(request.Username)

		// Send the prompt to the Inferer and receive the Response
		var Sender = inferer.Endpoint{
			Host: "localhost",
			Port: 7060,
		}
		response := Sender.SendMessage(convo, &w)

		// Append the response to Redis and MongoDB
		var assistantResponse = data.Message{
			Role:    "assistant",
			Content: response,
		}
		convo.Messages = append(convo.Messages, assistantResponse)
		rc.SetConversationData(request.Id.Hex(), convo)
		c.InsertNewMessage(request.Username, request.Id, assistantResponse)

		// Reply that everything is POG
		w.WriteHeader(http.StatusOK)
	}
}
