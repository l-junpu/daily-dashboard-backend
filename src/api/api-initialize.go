package api

import (
	"daily-dashboard-backend/src/database"
	"daily-dashboard-backend/src/llm"
	"net/http"
)

func InitializeMssqlApi(m *database.MssqlServer) {
	// User Related APIs
	http.HandleFunc("/login", HandleUserLogin(m))

	// Task Related APIs
	http.HandleFunc("/add_task_to_user", HandleAddTaskToUser(m))
	http.HandleFunc("/get_tasks_from_user", HandleGetTasksFromUser(m))
	http.HandleFunc("/update_task_for_user", HandleUpdateTaskForUser(m))
	http.HandleFunc("/remove_task_from_user", HandleRemoveTaskFromUser(m))
}

func InitializeMongoDBApi(c *llm.MongoDBClient, r *llm.RedisClient) {
	// Frontend Related APIs
	http.HandleFunc("/get_convos", HandleGetConvosFromUser(c))
	http.HandleFunc("/get_convo_details", HandleGetConvoDetails(c, r))
	http.HandleFunc("/get_convo_history", HandleGetConvoHistory(c, r))

	// LLM Related APIs
	http.HandleFunc("/create_new_convo", HandleCreateNewConvo(c))
	http.HandleFunc("/delete_convo", HandleDeleteConvo(c, r))
	http.HandleFunc("/new_user_prompt", HandleNewUserPrompt(c, r))
}

func InitializeSharedApi(m *database.MssqlServer, c *llm.MongoDBClient) {
	// User Related APIs
	http.HandleFunc("/register", HandleCommonRegistration(m, c))
}
