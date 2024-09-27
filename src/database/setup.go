package database

import (
	"database/sql"
	"fmt"
)

/*
Initialization Functions
*/
func (s *MssqlServer) createDashboarDatadDbIfNotExist() error {
	// Server=localhost\SQLEXPRESS;Database=master;Trusted_Connection=True;
	connectionString := fmt.Sprintf("server=%s; database=%s; trustedConnection=%t; trustServerCertificate=%t", s.ServerName, "master", s.TrustedConnection, s.TrustServerCertificate)
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer db.Close()

	createDatabaseQuery := `
	DECLARE @createDb NVARCHAR(MAX);
	IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = @databaseName)
	BEGIN
		SET @createDb = 'CREATE DATABASE [' + @databaseName + ']';
		EXEC sp_executesql @createDb;
	END;`

	stmt, err := db.Prepare(createDatabaseQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sql.Named("databaseName", s.DatabaseName))
	if err != nil {
		return fmt.Errorf("error executing named command: %w", err)
	}
	return nil
}

func (s *MssqlServer) createUserTableIfNotExist() error {
	createTableQuery := `
	IF OBJECT_ID('Users', 'U') IS NULL
	BEGIN
		CREATE TABLE Users (
			Id INT IDENTITY(1,1) PRIMARY KEY,
			Username VARCHAR(20) NOT NULL UNIQUE,
			Password VARCHAR(20) NOT NULL UNIQUE
		);
	END;`

	err := s.execCommand(createTableQuery)
	return err
}

func (s *MssqlServer) createTaskTableIfNotExist() error {
	createTableQuery := `
	IF OBJECT_ID('Tasks', 'U') IS NULL
	BEGIN
		CREATE TABLE Tasks (
		Id INT PRIMARY KEY IDENTITY(1,1),
		UserId INT NOT NULL,
		Text VARCHAR(2000) NOT NULL,
		LastModified DATETIME NOT NULL DEFAULT GETDATE(),
		CreatedOn DATETIME NOT NULL DEFAULT GETDATE(),
		CONSTRAINT FK_Tasks_Users FOREIGN KEY (UserId) REFERENCES Users(Id)
	);
	END;`

	err := s.execCommand(createTableQuery)
	return err
}