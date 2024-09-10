package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

/*
Database Operations
*/

func (s *MssqlServer) AddTaskToUser(username string, text string) error { // Consider changing username to userId
	userId, err := s.GetUserIdFromUsername(username)
	if err != nil {
		log.Fatal(err)
	}

	// If userId != -1, it means username exists
	if userId != -1 {
		addTask := `
		INSERT INTO Tasks (UserId, Text, LastModified, CreatedOn)
		VALUES (@userId, @text, GETDATE(), GETDATE());`

		err = s.execNamedCommand(addTask, sql.Named("userId", userId), sql.Named("text", text))
		return err
	}

	return err
}

func (s *MssqlServer) UpdateTaskForUser(taskId int, text string) error {
	updateTask := `
	UPDATE Tasks
	SET Text = @text, LastModified = GETDATE()
	WHERE Id = @taskId;`

	if err := s.execNamedCommand(updateTask, sql.Named("text", text), sql.Named("taskId", taskId)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
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

func (s *MssqlServer) GetWeeklyTasksFromUser(username string) error { // Consider changing username to userId
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
		friday := monday.AddDate(0, 0, 4)
		friday = time.Date(friday.Year(), friday.Month(), friday.Day(), 23, 59, 59, 0, time.UTC)

		getWeeklyTasks := `
		SELECT * FROM Tasks
		WHERE UserId = @userId AND @startDate <= CreatedOn AND @endDate >= CreatedOn;`

		weeklyTasks, err := s.execNamedQuery(getWeeklyTasks, sql.Named("userId", userId), sql.Named("startDate", monday), sql.Named("endDate", friday))
		if err != nil {
			return err
		}
		fmt.Println("Weekly Tasks:")
		s.printRows(weeklyTasks) // May want to consider writing a function to format it into an array of "struct/interface" so its easier to send to frontend
	}

	return nil
}

func (s *MssqlServer) FindTasksContainingText() {

}

func (s *MssqlServer) FindTasksContainingTextInPastDays() {

}
