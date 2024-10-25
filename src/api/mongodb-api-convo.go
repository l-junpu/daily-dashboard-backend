package api

import (
	"daily-dashboard-backend/src/data"
	"daily-dashboard-backend/src/llm"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

		WriteAsJson(w, map[string]interface{}{"id": id}, http.StatusOK)
	}
}
