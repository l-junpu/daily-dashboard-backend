package database

import (
	"daily-dashboard-backend/src/data"
	"database/sql"
	"fmt"
	"log"
	"time"
)

/*
Database Operations
*/

func (s *MssqlServer) AddTaskToUser(username string, title string, text string) (data.NewTaskResponse, error) { // Consider changing username to userId
	userId, err := s.GetUserIdFromUsername(username)
	if err != nil {
		log.Fatal(err)
	}

	// If userId != -1, it means username exists
	var response data.NewTaskResponse
	if userId != -1 {
		addTask := `
		INSERT INTO Tasks (UserId, Title, Text, Status, LastModified, CreatedOn)
		OUTPUT inserted.Id, inserted.LastModified, inserted.CreatedOn
		VALUES (@userId, @title, @text, 0, GETDATE(), GETDATE());`

		db, err := s.establishConnection()
		if err != nil {
			return response, err
		}
		defer db.Close()

		var response data.NewTaskResponse
		if err = db.QueryRow(addTask, sql.Named("userId", userId), sql.Named("title", title), sql.Named("text", text)).
			Scan(&response.TaskId, &response.LastModified, &response.CreatedOn); err != nil {
			return response, err
		}

		return response, nil
	}

	return response, err
}

func (s *MssqlServer) UpdateTaskForUser(taskId int, title string, text string, status bool) (string, error) {
	updateTask := `
	UPDATE Tasks
	SET Title = @title, Text = @text, Status = @status, LastModified = GETDATE()
	OUTPUT inserted.LastModified
	WHERE Id = @taskId;`

	db, err := s.establishConnection()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var lastModified string
	if err = db.QueryRow(updateTask, sql.Named("title", title), sql.Named("text", text), sql.Named("status", status), sql.Named("taskId", taskId)).
		Scan(&lastModified); err != nil {
		return "", err
	}
	
	return lastModified, nil
}

func (s *MssqlServer) RemoveTaskFromUser(taskId int) error {
	deleteTask := `
	DELETE FROM Tasks
	WHERE Id = @taskId;`

	if err := s.execNamedCommand(deleteTask, sql.Named("taskId", taskId)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *MssqlServer) GetWeeklyTasksFromUser(username string) ([]data.TaskDetailsResponse, error) {
	getMonday := func(now time.Time) time.Time {
		weekday := now.Weekday()
		if weekday == time.Sunday {
			return now.AddDate(0, 0, -6)
		} else if weekday == time.Saturday {
			return now.AddDate(0, 0, -5)
		} else {
			return now.AddDate(0, 0, -int(weekday)+1)
		}
	}

	userId, err := s.GetUserIdFromUsername(username)
	if err != nil {
		log.Fatal(err)
	}

	// If userId != -1, it means username exists
	if userId != -1 {
		monday := getMonday(time.Now())
		monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)
		sunday := monday.AddDate(0, 0, 6)
		sunday = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 23, 59, 59, 0, time.UTC)

		getWeeklyTasks := `
		SELECT Id, Title, Text, Status, LastModified, CreatedOn FROM Tasks
		WHERE UserId = @userId AND @startDate <= CreatedOn AND @endDate >= CreatedOn;`

		weeklyTaskRows, err := s.execNamedQuery(getWeeklyTasks, sql.Named("userId", userId), sql.Named("startDate", monday), sql.Named("endDate", sunday))
		if err != nil {
			return nil, err
		}
		defer weeklyTaskRows.Close()

		var weeklyTasks []data.TaskDetailsResponse
		for weeklyTaskRows.Next() {
			var t data.TaskDetailsResponse
			err = weeklyTaskRows.Scan(&t.TaskId, &t.Title, &t.Contents, &t.Status, &t.LastModified, &t.CreatedOn)
			if err != nil {
				log.Fatal(err)
			}
			weeklyTasks = append(weeklyTasks, t)
		}

		return weeklyTasks, nil
	}

	return nil, nil
}
