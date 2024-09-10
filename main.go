package main

import (
	"daily-dashboard-backend/src/database"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	svr := database.MssqlServer{
		ServerName:             "localhost\\SQLEXPRESS",
		DatabaseName:           "DashboardData",
		TrustedConnection:      true,
		TrustServerCertificate: true,
		EnablePrintouts:        false,
	}

	if err := svr.Initialise(); err != nil {
		log.Fatal(err)
	}

	svr.RegisterNewUser("sussy user")
	//svr.AddTaskToUser("sussy user", "i love my terrapin")
	//svr.UpdateTaskForUser(2, "booooo")
	// svr.RemoveTaskFromUser(10)
	svr.GetWeeklyTasksFromUser("sussy user")

	svr.PrintDebugData()
}
