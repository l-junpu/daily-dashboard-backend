package api

import (
	"daily-dashboard-backend/src/database"
	"net/http"
)

func InitializeApi(s *database.MssqlServer) {
	// User Related APIs
	http.HandleFunc("/register", HandleRegisterNewUser(s))
	http.HandleFunc("/get_user_id", HandleGetUserId(s))

	// Task Related APIs
	http.HandleFunc("/add_task_to_user", HandleAddTaskToUser(s))
	http.HandleFunc("/update_task_for_user", HandleUpdateTaskForUser(s))
	http.HandleFunc("/remove_task_from_user", HandleRemoveTaskFromUser(s))
}
