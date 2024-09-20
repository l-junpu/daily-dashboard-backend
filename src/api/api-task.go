package api

import (
	"daily-dashboard-backend/src/data"
	"daily-dashboard-backend/src/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HandleAddTaskToUser(s *database.MssqlServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)

		var task data.NewTask
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			WriteAsJson(w, map[string]interface{}{"status": false})
			err = fmt.Errorf("unable to decode task data in HandleAddTaskToUser(): %w", err)
			log.Fatal(err)
		}

		if err := s.AddTaskToUser(task.Username, task.Contents); err != nil {
			WriteAsJson(w, map[string]interface{}{"status": false})
			err = fmt.Errorf("unable to add task to user in HandleAddTaskToUser(): %w", err)
			log.Fatal(err)
		}

		WriteAsJson(w, map[string]interface{}{"status": true})
	}
}

func HandleUpdateTaskForUser(s *database.MssqlServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)

		var task data.UpdateTask
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			WriteAsJson(w, map[string]interface{}{"status": false})
			err = fmt.Errorf("unable to decode task data in HandleUpdateTaskForUser(): %w", err)
			log.Fatal(err)
		}

		if err := s.UpdateTaskForUser(task.TaskId, task.Contents); err != nil {
			WriteAsJson(w, map[string]interface{}{"status": false})
			err = fmt.Errorf("unable to update task for user in HandleUpdateTaskForUser(): %w", err)
			log.Fatal(err)
		}

		WriteAsJson(w, map[string]interface{}{"status": true})
	}
}

func HandleRemoveTaskFromUser(s *database.MssqlServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)

		var task data.TaskId
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			WriteAsJson(w, map[string]interface{}{"status": false})
			err = fmt.Errorf("unable to decode task data in HandleRemoveTaskFromUser(): %w", err)
			log.Fatal(err)
		}

		if err := s.RemoveTaskFromUser(task.TaskId); err != nil {
			WriteAsJson(w, map[string]interface{}{"status": false})
			err = fmt.Errorf("unable to remove task from user in HandleRemoveTaskFromUser(): %w", err)
			log.Fatal(err)
		}

		WriteAsJson(w, map[string]interface{}{"status": true})
	}
}
