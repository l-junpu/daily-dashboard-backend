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
		Debugging Functions
	*/
	PrintDebugData()
	getLastFewUsers()
	getLastFewTasks()
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
