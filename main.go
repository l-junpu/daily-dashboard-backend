package main

import (
	"daily-dashboard-backend/src/api"
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

	// Initialize MSSQL Server
	if err := svr.Initialise(); err != nil {
		log.Fatal(err)
	}

	// Register Http Handler Functions
	api.InitializeApi(&svr)

	svr.Tick()
}
