package main

import (
	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	//mongoDbUri := "mongodb://localhost:27017/"

	// sqlSvr := database.MssqlServer{
	// 	ServerName:             "localhost\\SQLEXPRESS",
	// 	DatabaseName:           "DashboardData",
	// 	TrustedConnection:      true,
	// 	TrustServerCertificate: true,
	// 	EnablePrintouts:        false,
	// }

	// // Initialize MSSQL Server
	// if err := sqlSvr.Initialise(); err != nil {
	// 	log.Fatal(err)
	// }

	// // Register Http Handler Functions For
	// api.InitializeApi(&sqlSvr)

	// http.ListenAndServe("localhost:8080", nil)

	// client, err := llm.CreateMongoDBClient(mongoDbUri)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
