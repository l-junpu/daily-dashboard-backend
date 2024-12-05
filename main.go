package main

import (
	"daily-dashboard-backend/src/api"
	"daily-dashboard-backend/src/database"
	"daily-dashboard-backend/src/inferer"
	"daily-dashboard-backend/src/llm"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	// Create & Initialize MSSQL Server
	sqlSvr := database.CreateMssqlServer()
	if err := sqlSvr.Initialise(); err != nil {
		log.Fatal(err)
	}
	api.InitializeMssqlApi(sqlSvr)

	// Create Redis Server
	redisClient, err := llm.CreateRedisClient()
	if err != nil {
		log.Fatal(err)
	}

	// Create an Endpoint Scheduler
	scheduler, err := inferer.CreateScheduler(os.Getwd())
	if err != nil {
		log.Fatal(err)
	}

	// Create & Initialize MongoDB Server + Redis Client
	mongoClient, err := llm.CreateMongoDBClient()
	if err != nil {
		log.Fatal(err)
	}
	api.InitializeMongoDBApi(mongoClient, redisClient, scheduler)

	// Register Common APIs - MSSQL + MongoDB
	api.InitializeSharedApi(sqlSvr, mongoClient)

	// Listen to requests to Server
	serverAddr := os.Getenv("SERVER_ADDRESS")
	http.ListenAndServe(serverAddr, nil)
}
