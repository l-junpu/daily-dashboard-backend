package api

import (
	"daily-dashboard-backend/src/data"
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

func HandleGetConvoHistory(c *llm.MongoDBClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}

		// Decode user data
		var request data.ConversationRequest
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

func HandleDeleteConvo(c *llm.MongoDBClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}

		// Decode user data
		var request data.ConversationRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("unable to decode delete convo details in HandleDeleteConvo(): %w", err)
			log.Fatal(err)
		}

		if err := c.DeleteConversation(request.Username, request.ObjectID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("unable to delete conversation from MongoDB: %w", err)
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func HandleNewUserPrompt(c *llm.MongoDBClient) http.HandlerFunc {
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

		// establish stream to frontend - todo

		// reply first
		w.WriteHeader(http.StatusOK)
	}
}
