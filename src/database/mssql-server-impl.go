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
	if err = s.createWeeklyTasksTableIfNotExist(db); err != nil {
		return err
	}
	if err = s.createDailyTasksTableIfNotExist(db); err != nil {
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
	connectionString := fmt.Sprintf("server=%s; database=%s; trustedConnection=True; trustServerCertificate=%t", s.serverName, "master", s.trustServerCertificate)
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	createDatabaseQuery := `
	IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = '%s')
	BEGIN
		CREATE DATABASE [%s];
	END
	`
	createDatabaseQuery = fmt.Sprintf(createDatabaseQuery, s.databaseName, s.databaseName)

	_, err = db.Exec(createDatabaseQuery)
	if err != nil {
		return fmt.Errorf("error creating DashboardData database: %w", err)
	}

	fmt.Println("DashboardData database created or already exists!")
	return nil
}

func (s *MssqlServer) createUserTableIfNotExist(db *sql.DB) error {
	createTableQuery := `
	IF OBJECT_ID('Users', 'U') IS NULL
	BEGIN
		CREATE TABLE Users (
			Id INT IDENTITY(1,1) PRIMARY KEY,
			Username VARCHAR(20) NOT NULL,
			Password VARCHAR(20) NOT NULL
		);
	END;`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating Users table: %w", err)
	}

	fmt.Println("Users table created or already exists!")
	return nil
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
		Day DATE NOT NULL,
		CONSTRAINT FK_Tasks_Users FOREIGN KEY (UserId) REFERENCES Users(Id)
	);
	END;`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating Tasks table: %w", err)
	}

	fmt.Println("Tasks table created or already exists!")
	return nil
}

func (s *MssqlServer) createDailyTasksTableIfNotExist(db *sql.DB) error {
	createTableQuery := `
	IF OBJECT_ID('DailyTasks', 'U') IS NULL
	BEGIN
		CREATE TABLE DailyTasks (
			Id INT IDENTITY(1,1) PRIMARY KEY,
			WeeklyTaskId INT NOT NULL,
			TaskId INT NOT NULL,
			Day DATE NOT NULL,
			CONSTRAINT FK_DailyTasks_WeeklyTasks FOREIGN KEY (WeeklyTaskId) REFERENCES WeeklyTasks(Id),
			CONSTRAINT FK_DailyTasks_Tasks FOREIGN KEY (TaskId) REFERENCES Tasks(Id)
		);
	END;`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating Daily Tasks table: %w", err)
	}

	fmt.Println("Daily Tasks table created or already exists!")
	return nil
}

func (s *MssqlServer) createWeeklyTasksTableIfNotExist(db *sql.DB) error {
	createTableQuery := `
	IF OBJECT_ID('WeeklyTasks', 'U') IS NULL
	BEGIN
		CREATE TABLE WeeklyTasks (
		Id INT PRIMARY KEY IDENTITY(1,1),
		UserId INT NOT NULL,
		Week INT NOT NULL,
		Year INT NOT NULL,
		CONSTRAINT FK_WeeklyTasks_Users FOREIGN KEY (UserId) REFERENCES Users(Id)
	);
	END;`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating Weekly Tasks table: %w", err)
	}

	fmt.Println("Weekly Tasks table created or already exists!")
	return nil
}

func (s *MssqlServer) PrintDebugData() {
	printRows := func(category string, rows *sql.Rows) {
		defer rows.Close()
		category = fmt.Sprintf("Displaying '%s':", category)

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

	// Do it for the rest
}

func (s *MssqlServer) getLastFewUsers() (*sql.Rows, error) {
	return nil, nil
}

func (s *MssqlServer) getLastFewTasks() (*sql.Rows, error) {
	return nil, nil
}

func (s *MssqlServer) getLastFewDailyTasks() (*sql.Rows, error) {
	return nil, nil
}

func (s *MssqlServer) getLastFewWeeklyTasks() (*sql.Rows, error) {
	return nil, nil
}
