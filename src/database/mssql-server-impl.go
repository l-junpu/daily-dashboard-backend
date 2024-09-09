package database

import (
	"database/sql"
	"fmt"
	"log"
)

/*
Generic Connection Functions
*/
func (s *MssqlServer) Initialise() error {
	if err := s.createDashboarDatadDbIfNotExist(); err != nil {
		return err
	}

	db, err := s.establishConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	if err = s.createUserTableIfNotExist(db); err != nil {
		return err
	}
	if err = s.createTaskTableIfNotExist(db); err != nil {
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
	connectionString := fmt.Sprintf("server=%s; database=%s; trustedConnection=True; trustServerCertificate=%t", s.serverName, s.databaseName, s.trustServerCertificate)
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

	addDataQuery := `
	INSERT INTO Users (Username)
	VALUES (@username);`

	stmt, err := db.Prepare(addDataQuery)
	if err != nil {
		print(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sql.Named("username", username))
	if err != nil {
		print("err 2: ", err)
		return err
	}

	// //addDataQuery = fmt.Sprintf(addDataQuery, username)
	// _, err = db.Exec(addDataQuery)
	// if err != nil {
	// 	return fmt.Errorf("error executing command: %w", err)
	// }
	// defer rows.Close()
	// // Get the column names
	// columns, err := rows.Columns()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Print the column names
	// fmt.Println(columns)

	// // Iterate over the rows
	// for rows.Next() {
	// 	values := make([]*string, len(columns))
	// 	for i := range values {
	// 		values[i] = new(string)
	// 		if err = rows.Scan(values[i]); err != nil {
	// 			return fmt.Errorf("error scanning for rows: %w", err)
	// 		}
	// 		fmt.Println(*(values[i]))
	// 	}
	// }

	// // Check for any errors during iteration
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }

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
	connectionString := fmt.Sprintf("server=%s; database=%s; trustedConnection=True; trustServerCertificate=%t", s.serverName, "master", s.trustServerCertificate)
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	createDatabaseQuery := `
	IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = '%s')
	BEGIN
		CREATE DATABASE [%s];
	END;`
	createDatabaseQuery = fmt.Sprintf(createDatabaseQuery, s.databaseName, s.databaseName)

	_, err = db.Exec(createDatabaseQuery)
	if err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}
	return err
}

func (s *MssqlServer) createUserTableIfNotExist(db *sql.DB) error {
	createTableQuery := `
	IF OBJECT_ID('Users', 'U') IS NULL
	BEGIN
		CREATE TABLE Users (
			Id INT IDENTITY(1,1) PRIMARY KEY,
			Username VARCHAR(20) NOT NULL
		);
	END;`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}
	return err
}

func (s *MssqlServer) createTaskTableIfNotExist(db *sql.DB) error {
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
	
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}
	return err
}

func (s *MssqlServer) PrintDebugData() {
	printRows := func(category string, rows *sql.Rows) {
		defer rows.Close()
		category = fmt.Sprintf("Displaying '%s':", category)
		fmt.Println(category)

		for rows.Next() {
			var data string
			err := rows.Scan(&data)
			if err != nil {
				log.Fatal("error scanning row: ", err)
			}
			fmt.Println(data)
		}
	}

	userData, err := s.getLastFewUsers()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	printRows("User Data", userData)

	taskData, err := s.getLastFewTasks()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	printRows("Task Data", taskData)

	// Do it for the rest
}

func (s *MssqlServer) getLastFewUsers() (*sql.Rows, error) {
	return nil, nil
}

func (s *MssqlServer) getLastFewTasks() (*sql.Rows, error) {
	return nil, nil
}
