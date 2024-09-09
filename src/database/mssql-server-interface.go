package database

type MssqlServerInterface interface {
	/*
		Generic Connection Functions
	*/
	Initialise()
	establishConnection()
	generateConnectionString()

	/*
		User Initialization
	*/
	RegisterNewUser()
	GetUserIdFromUsername()

	/*
		Database Operations
	*/
	AddUserToDb()
	AddTaskToUser()
	RemoveTaskFromUser()
	GetWeeklyTasksFromUser()
	FindTasksContainingText()
	FindTasksContainingTextInPastDays()

	/*
		Initialization Functions
	*/
	createDashboarDatadDbIfNotExist()
	createUserTableIfNotExist()
	createTaskTableIfNotExist()

	/*
		Exec SQL Commands
	*/
	execQuery()
	execCommand()
	execNamedCommand()
	printRows()

	/*
		Debugging Functions
	*/
	PrintDebugData()
	getLastFewUsers()
	getLastFewTasks()
}

type MssqlServer struct {
	ServerName             string
	DatabaseName           string
	TrustedConnection      bool
	TrustServerCertificate bool
}
