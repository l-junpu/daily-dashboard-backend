package database

import (
	"database/sql"
	"log"
)

/*
Database Operations
*/

func (s *MssqlServer) AddTaskToUser(username string, text string) error {
	userId, err := s.GetUserIdFromUsername(username)
	if err != nil {
		log.Fatal(err)
	}

	if userId != -1 {
		addTask := `
		INSERT INTO Tasks (UserId, Text, LastModified, CreatedOn)
		VALUES (@userId, @text, GETDATE(), GETDATE());`
	
		err = s.execNamedCommand(addTask, sql.Named("userId", userId), sql.Named("text", text))
		return err
	}

	return err
}

func (s *MssqlServer) UpdateTaskForUser() {

}

func (s *MssqlServer) RemoveTaskFromUser() {

}

func (s *MssqlServer) GetWeeklyTasksFromUser() {

}

func (s *MssqlServer) FindTasksContainingText() {

}

func (s *MssqlServer) FindTasksContainingTextInPastDays() {

}
