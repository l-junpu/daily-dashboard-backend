package api

import (
	"daily-dashboard-backend/src/database"
	"daily-dashboard-backend/src/llm"
	"net/http"
)

func InitializeMssqlApi(m *database.MssqlServer) {
	// User Related APIs
	http.HandleFunc("/register", HandleRegisterNewUser(m)) // Shift this to SharedApi
	http.HandleFunc("/login", HandleUserLogin(m))

	// Task Related APIs
	http.HandleFunc("/add_task_to_user", HandleAddTaskToUser(m))
	http.HandleFunc("/get_tasks_from_user", HandleGetTasksFromUser(m))
	http.HandleFunc("/update_task_for_user", HandleUpdateTaskForUser(m))
	http.HandleFunc("/remove_task_from_user", HandleRemoveTaskFromUser(m))
}

func InitializeMongoDBApi(c *llm.MongoDBClient) {
	/*
		// User Related APIs
		http.HandleFunc("/register", HandleRegisterNewUser(c)) // Shift this to SharedApi

		// Frontend Related APIs
		http.HandleFunc("/get_convos", HandleGetConvosFromUser(c))
		http.HandleFunc("/get_convo_details", HandleGetConvoDetails(c))

		// LLM Related APIs
		http.HandleFunc("/create_new_convo", HandleCreateNewConvo(s))
		http.HandleFunc("/delete_convo", HandleDeleteConvo(s))
		http.HandleFunc("/new_user_prompt", HandleNewUserPrompt(s))

	*/
}

func InitializeSharedApi(m *database.MssqlServer, c *llm.MongoDBClient) {

}
