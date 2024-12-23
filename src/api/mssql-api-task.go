package api

import (
	"daily-dashboard-backend/src/data"
	"daily-dashboard-backend/src/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HandleGetTasksFromUser(s *database.MssqlServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}

		var user data.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("unable to decode user data in HandleGetTasksFromUser(): %w", err)
			log.Fatal(err)
		}

		fmt.Println(user)

		tasks, err := s.GetWeeklyTasksFromUser(user.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("unable to add task to user in HandleAddTaskToUser(): %w", err)
			log.Fatal(err)
		}

		WriteAsJson(w, map[string]interface{}{"tasks": tasks}, http.StatusOK)
	}
}

func HandleAddTaskToUser(s *database.MssqlServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}

		var task data.NewTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("unable to decode task data in HandleAddTaskToUser(): %w", err)
			log.Fatal(err)
		}

		fmt.Println(task)

		response, err := s.AddTaskToUser(task.Username, task.Title, task.Contents)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("unable to add task to user in HandleAddTaskToUser(): %w", err)
			log.Fatal(err)
		}

		WriteAsJson(w, map[string]interface{}{"taskId": response.TaskId, "lastModified": response.LastModified, "createdOn": response.CreatedOn}, http.StatusOK)
	}
}

func HandleUpdateTaskForUser(s *database.MssqlServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}

		var task data.UpdateTaskContentsRequest
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("unable to decode task data in HandleUpdateTaskForUser(): %w", err)
			log.Fatal(err)
		}

		lastModified, err := s.UpdateTaskForUser(task.TaskId, task.Title, task.Contents, task.Status);
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("unable to update task for user in HandleUpdateTaskForUser(): %w", err)
			log.Fatal(err)
		}

		// Change this to return LastModified time
		WriteAsJson(w, map[string]interface{}{"lastModified": lastModified}, http.StatusOK)
	}
}

func HandleRemoveTaskFromUser(s *database.MssqlServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if preflight := HandleOptionsPreflightRequests(w, r); preflight {
			return
		}

		var task data.RemoveTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = fmt.Errorf("unable to decode task data in HandleRemoveTaskFromUser(): %w", err)
			log.Fatal(err)
		}

		if err := s.RemoveTaskFromUser(task.TaskId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("unable to remove task from user in HandleRemoveTaskFromUser(): %w", err)
			log.Fatal(err)
		}

		WriteAsJson(w, map[string]interface{}{"status": true}, http.StatusOK)
	}
}
