package database

import "database/sql"

type MssqlServerInterface interface {
	/*
		Generic Connection Functions
	*/
	EstablishConnection()
	TerminateConnection()
	generateConnectionString()

	/*
		User Initialization
	*/
	RegisterNewUser()
	GetUserIdFromUsername()

	/*
		Database Operations
	*/
	AddTaskToUser()
	RemoveTaskFromUser()
	GetWeeklyTasksFromUser()
	FindTasksContainingText()
	FindTasksContainingTextInPastDays()

	/*
		Initialization Functions
	*/
	createUserTableIfNotExist()
	createTaskTableIfNotExist()
	createDailyTasksTableIfNotExist()
	createWeeklyTasksTableIfNotExist()
}

type MssqlServer struct {
	serverName             string
	databaseName           string
	userId                 string
	password               string
	trustServerCertificate bool

	db *sql.DB
}

func CreateMssqlServer(serverName string, databaseName string, userId string, password string, trustServerCertificate bool) *MssqlServer {
	s := MssqlServer{
		serverName:             serverName,
		databaseName:           databaseName,
		userId:                 userId,
		password:               password,
		trustServerCertificate: trustServerCertificate,
		db:                     nil,
	}

	return &s
}
