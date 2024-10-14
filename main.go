package main

import (
	"daily-dashboard-backend/src/api"
	"daily-dashboard-backend/src/database"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Getenv("MSSQL_USERNAME")
	password := os.Getenv("MSSQL_PASSWORD")

	sqlSvr := database.MssqlServer{
		ServerName:             "localhost",
		Port:                   1433,
		Username:               username,
		Password:               password,
		DatabaseName:           "DashboardData",
		TrustedConnection:      true,
		TrustServerCertificate: true,
		EnablePrintouts:        false,
	}

	// Initialize MSSQL Server
	if err := sqlSvr.Initialise(); err != nil {
		log.Fatal(err)
	}

	// Register Http Handler Functions For
	api.InitializeApi(&sqlSvr)

	http.ListenAndServe("localhost:8080", nil)
}
