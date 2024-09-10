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
	}

	if err := svr.Initialise(); err != nil {
		log.Fatal(err)
	}

	svr.RegisterNewUser("jun pu")
	svr.AddTaskToUser("ryanchoo", "i love my terrapin")

	// if err := svr.AddUserToDb("sample user"); err != nil {
	// 	log.Fatal(err)
	// }

	svr.PrintDebugData()
}
