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
	createDailyTasksTableIfNotExist()
	createWeeklyTasksTableIfNotExist()

	/*
		Debugging Functions
	*/
	PrintDebugData()
	getLastFewUsers()
	getLastFewTasks()
	getLastFewDailyTasks()
	getLastFewWeeklyTasks()
}

type MssqlServer struct {
	serverName             string
	databaseName           string
	trustServerCertificate bool
}

func CreateMssqlServer(serverName string, databaseName string, trustServerCertificate bool) *MssqlServer {
	s := MssqlServer{
		serverName:             serverName,
		databaseName:           databaseName,
		trustServerCertificate: trustServerCertificate,
	}

	return &s
}
