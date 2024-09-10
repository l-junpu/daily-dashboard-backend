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
	AddTaskToUser()
	UpdateTaskForUser()
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
	execNamedQuery()
	printRows()

	/*
		Debugging Functions
	*/
	PrintDebugData()
	getLastFewUsers()
	getLastFewTasks()
}
