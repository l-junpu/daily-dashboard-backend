package database

import (
	"database/sql"
	"fmt"
)

/*
User Initialization
*/
func (s *MssqlServer) RegisterNewUser(username string) error {
	db, err := s.establishConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	addUserCommand := `
	IF NOT EXISTS (SELECT Username from Users WHERE Username = @username)
	BEGIN
		INSERT INTO Users (Username)
		VALUES (@username);
	END;`

	err = s.execNamedCommand(addUserCommand, sql.Named("username", username))
	return err
}

func (s *MssqlServer) GetUserIdFromUsername(username string) (int, error) {
	db, err := s.establishConnection()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	// Retrieve primary key
	queryUserId := `
	SELECT Id from Users
	WHERE Username = @username`

	rows, err := s.execNamedQuery(queryUserId, sql.Named("username", username))
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	// Working with the "guarantee" that there are only
	// unique Usernames
	var userId int
	rows.Next()
	if err = rows.Scan(&userId); err != nil {
		return -1, fmt.Errorf("unable to scan for userId: %w", err)
	}

	return userId, nil
}