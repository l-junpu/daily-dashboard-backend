package database

import (
	"database/sql"
	"fmt"
)

// RegisterNewUser creates a new user in the database if the username does not already exist.
// @param username The username to register.
// @return An error if the registration fails, otherwise nil.
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

// GetUserIdFromUsername retrieves the user ID associated with the given username.
// @param username The username to retrieve the ID for.
// @return The user ID if found, -1 if not found, and an error if the retrieval fails.
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

	// Check if there are any results
	// If there are no results, return -1 and not an error
	if !rows.Next() {
		return -1, nil
	}

	// Working with the "guarantee" that there are only
	// unique Usernames
	var userId int
	if err = rows.Scan(&userId); err != nil {
		return -1, fmt.Errorf("unable to scan for userId: %w", err)
	}

	return userId, nil
}
