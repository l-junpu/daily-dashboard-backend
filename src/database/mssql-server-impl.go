package database

import (
	"database/sql"
	"fmt"
)

/*
Generic Connection Functions
*/
func (s *MssqlServer) EstablishConnection() error {
	connectionString := s.generateConnectionString()
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	s.db = db
	if err = db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	s.createUserTableIfNotExist()
	s.createTaskTableIfNotExist()
	s.createDailyTasksTableIfNotExist()
	s.createWeeklyTasksTableIfNotExist()

	return nil
}

func (s *MssqlServer) TerminateConnection() {
	if err := s.db.Close(); err != nil {
		fmt.Println(err)
	}
}

func (s *MssqlServer) generateConnectionString() string {
	connectionString := fmt.Sprintf("server=%s; user id=%s; password=%s; database=%s; trustServerCertificate=%t", s.serverName, s.userId, s.password, s.databaseName, s.trustServerCertificate)
	return connectionString
}

/*
User Initialization
*/
func (s *MssqlServer) RegisterNewUser(username string, password string) error {
	return fmt.Errorf("")
}

func (s *MssqlServer) GetUserIdFromUsername(username string) error {
	return fmt.Errorf("")
}

/*
Database Operations
*/
func (s *MssqlServer) AddTaskToUser() {

}

func (s *MssqlServer) RemoveTaskFromUser() {

}

func (s *MssqlServer) GetWeeklyTasksFromUser() {

}

func (s *MssqlServer) FindTasksContainingText() {

}

func (s *MssqlServer) FindTasksContainingTextInPastDays() {

}

/*
Initialization Functions
*/
func (s *MssqlServer) createUserTableIfNotExist() {

}

func (s *MssqlServer) createTaskTableIfNotExist() {

}

func (s *MssqlServer) createDailyTasksTableIfNotExist() {

}

func (s *MssqlServer) createWeeklyTasksTableIfNotExist() {

}
