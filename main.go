package main

import (
	"daily-dashboard-backend/src/database"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	svr := database.CreateMssqlServer("localhost", "", "", "master", true)
	if err := svr.EstablishConnection(); err != nil {
		fmt.Println(err)
	}
	defer svr.TerminateConnection()
}
