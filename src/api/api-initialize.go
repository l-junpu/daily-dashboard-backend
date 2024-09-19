package api

import (
	"daily-dashboard-backend/src/database"
	"net/http"
)

func InitializeApi(s *database.MssqlServer) {
	http.HandleFunc("/register", HandleRegisterNewUser(s))
	http.HandleFunc("/get_user_id", HandleGetUserId(s))
}
