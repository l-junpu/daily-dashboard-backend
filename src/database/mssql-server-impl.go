package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

/*
Generic Connection Functions
*/
func (s *MssqlServer) Initialise() error {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)

	if err := s.createDashboarDatadDbIfNotExist(); err != nil {
		return err
	}

	if err := s.createUserTableIfNotExist(); err != nil {
		return err
	}
	if err := s.createTaskTableIfNotExist(); err != nil {
		return err
	}

	return nil
}

func (s *MssqlServer) establishConnection() (*sql.DB, error) {
	connectionString := s.generateConnectionString()
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return db, nil
}

func (s *MssqlServer) generateConnectionString() string {
	connectionString := fmt.Sprintf("server=%s; database=%s; trustedConnection=%t; trustServerCertificate=%t", s.ServerName, s.DatabaseName, s.TrustedConnection, s.TrustServerCertificate)
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
func (s *MssqlServer) AddUserToDb(username string) error {
	db, err := s.establishConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	addUserCommand := `
	INSERT INTO Users (Username)
	VALUES (@username);`

	err = s.execNamedCommand(addUserCommand, sql.Named("username", username))
	return err
}

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
func (s *MssqlServer) createDashboarDatadDbIfNotExist() error {
	// Server=localhost\SQLEXPRESS;Database=master;Trusted_Connection=True;
	connectionString := fmt.Sprintf("server=%s; database=%s; trustedConnection=%t; trustServerCertificate=%t", s.ServerName, "master", s.TrustedConnection, s.TrustServerCertificate)
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer db.Close()

	createDatabaseQuery := `
	IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = @databaseName)
	BEGIN
		CREATE DATABASE @databaseName;
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
			Username VARCHAR(20) NOT NULL
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

func (s *MssqlServer) execQuery(query string) (*sql.Rows, error) {
	db, err := s.establishConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	return rows, nil
}

func (s *MssqlServer) execCommand(command string) error {
	db, err := s.establishConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(command)
	if err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}
	return nil
}

func (s *MssqlServer) execNamedCommand(namedCommand string, namedArgs ...interface{}) error {
	db, err := s.establishConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(namedCommand)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedArgs...)
	if err != nil {
		return fmt.Errorf("error executing named command: %w", err)
	}
	return nil
}

func (s *MssqlServer) printRows(rows *sql.Rows) {
	// Print Headers
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(columns)

	// Prepare a slice of interfaces to hold the values for each column
	values := make([]interface{}, len(columns))
	valuePointers := make([]interface{}, len(columns))

	for rows.Next() {
		for i := range values {
			valuePointers[i] = &values[i]
		}

		if err = rows.Scan(valuePointers...); err != nil {
			log.Fatal(fmt.Errorf("error scanning for rows: %w", err))
		}

		for _, val := range values {
			var v string
			// Check for byte array - Strings
			// Else, any other integral type
			if b, ok := val.([]byte); ok {
				v = string(b)
			} else {
				v = fmt.Sprintf("%v", val)
			}
			fmt.Print(v, ", ")
		}
		fmt.Println()
	}
	fmt.Println()
}

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
