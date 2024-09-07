package main

import (
	"daily-dashboard-backend/src/database"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	svr := database.CreateMssqlServer("localhost\\SQLEXPRESS", "DashboardData", true)
	if err := svr.Initialise(); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}
