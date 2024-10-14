package api

import (
	"daily-dashboard-backend/src/data"
	"daily-dashboard-backend/src/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HandleRegisterNewUser(s *database.MssqlServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}
		
		// Decode user data
		var user data.UserDetails
		fmt.Println(r.Body)
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			WriteAsJson(w, map[string]interface{}{"status": false})
			err = fmt.Errorf("unable to decode user data in HandleRegisterNewUser(): %w", err)
			log.Fatal(err)
		}

		// Add username to database if it doesn't exist
		if err := s.RegisterNewUser(user.Username, user.Password); err != nil {
			WriteAsJson(w, map[string]interface{}{"status": false})
			err = fmt.Errorf("unable to insert username into database in HandleRegisterNewUser(): %w", err)
			log.Fatal(err)
		}

		WriteAsJson(w, map[string]interface{}{"status": true})
	}
}

func HandleUserLogin(s *database.MssqlServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}
		
		// Decode user data
		var user data.UserDetails
		fmt.Println(r.Body)
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			WriteAsJson(w, map[string]interface{}{"status": false})
			err = fmt.Errorf("unable to decode user data in HandleUserLogin(): %w", err)
			log.Fatal(err)
		}
		
		// Add username to database if it doesn't exist
		isValidLogin, err := s.VerifyUserLogin(user.Username, user.Password);
		fmt.Println("Verified user login request")
		if err != nil {
			WriteAsJson(w, map[string]interface{}{"status": false})
			err = fmt.Errorf("unable to verify user login details in HandleUserLogin(): %w", err)
			log.Fatal(err)
		}
		if !isValidLogin {
			WriteAsJson(w, map[string]interface{}{"status": false})
			return
		}

		WriteAsJson(w, map[string]interface{}{"status": true})
	}
}