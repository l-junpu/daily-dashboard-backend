package api

import (
	"daily-dashboard-backend/src/data"
	"daily-dashboard-backend/src/database"
	"daily-dashboard-backend/src/llm"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HandleCommonRegistration(m *database.MssqlServer, c *llm.MongoDBClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}

		// Decode user data
		var user data.UserDetails
		fmt.Println(r.Body)
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("unable to decode user data in HandleCommonRegistration(): %w", err)
			log.Fatal(err)
		}

		// Add username to MSSQL Database if it doesn't exist
		if err := m.RegisterNewUser(user.Username, user.Password); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("unable to insert username into MSSQL in HandleCommonRegistration(): %w", err)
			log.Fatal(err)
		}

		// Add username to MongoDB if it doesn't exist
		if err := c.InsertUser(user.Username, user.Password); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("unable to insert username into MongoDB in HandleCommonRegistration(): %w", err)
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
	}
}
