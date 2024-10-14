package database

import (
	"database/sql"
	"log"
)

/*
Display Debug Information from User / Task Table
*/
func (s *MssqlServer) PrintDebugData() {
	userData, err := s.getLastFewUsers()
	if err != nil {
		log.Fatal(err)
	}
	defer userData.Close()
	s.printRows(userData)

	taskData, err := s.getLastFewTasks()
	if err != nil {
		log.Fatal(err)
	}
	defer taskData.Close()
	s.printRows(taskData)
}

func (s *MssqlServer) getLastFewUsers() (*sql.Rows, error) {
	selectLastTenUsers := `
	SELECT TOP 10 *
	FROM Users
	ORDER BY Id DESC;`

	rows, err := s.execQuery(selectLastTenUsers)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *MssqlServer) getLastFewTasks() (*sql.Rows, error) {
	selectLastTenTasks := `
	SELECT TOP 10 *
	FROM Tasks
	ORDER BY Id DESC;`

	rows, err := s.execQuery(selectLastTenTasks)
	if err != nil {
		return nil, err
	}
	return rows, nil
}