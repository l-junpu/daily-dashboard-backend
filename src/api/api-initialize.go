package api

import (
	"daily-dashboard-backend/src/database"
	"net/http"
)

func InitializeApi(s *database.MssqlServer) {
	// User Related APIs
	http.HandleFunc("/register", HandleRegisterNewUser(s))
	http.HandleFunc("/login", HandleUserLogin(s))

	// Task Related APIs
	http.HandleFunc("/add_task_to_user", HandleAddTaskToUser(s))
	http.HandleFunc("/get_tasks_from_user", HandleGetTasksFromUser(s))
	http.HandleFunc("/update_task_for_user", HandleUpdateTaskForUser(s))
	http.HandleFunc("/remove_task_from_user", HandleRemoveTaskFromUser(s))
}
