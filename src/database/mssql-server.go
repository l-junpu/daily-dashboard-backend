package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type MssqlServer struct {
	ServerName             string
	Port                   int
	Username               string
	Password               string
	DatabaseName           string
	TrustedConnection      bool
	TrustServerCertificate bool
	EnablePrintouts        bool
}

/*
Connection Functions
*/
func CreateMssqlServer() *MssqlServer {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Port, err := strconv.Atoi(os.Getenv("MSSQL_PORT"))
	if err != nil {
		log.Fatal("Invalid port number in .env file")
	}

	TrustedConnection, err := strconv.ParseBool(os.Getenv("MSSQL_TRUSTED_CONNECTION"))
	if err != nil {
		log.Fatal("Invalid bool value in .env file")
	}

	TrustServerCertificate, err := strconv.ParseBool(os.Getenv("MSSQL_TRUST_SERVER_CERTIFICATE"))
	if err != nil {
		log.Fatal("Invalid bool value in .env file")
	}

	EnablePrintouts, err := strconv.ParseBool(os.Getenv("MSSQL_ENABLE_PRINTOUTS"))
	if err != nil {
		log.Fatal("Invalid bool value in .env file")
	}

	sqlSvr := MssqlServer{
		ServerName:             os.Getenv("MSSQL_SERVER_NAME"),
		Port:                   Port,
		Username:               os.Getenv("MSSQL_USERNAME"),
		Password:               os.Getenv("MSSQL_PASSWORD"),
		DatabaseName:           os.Getenv("MSSQL_DATABASE_NAME"),
		TrustedConnection:      TrustedConnection,
		TrustServerCertificate: TrustServerCertificate,
		EnablePrintouts:        EnablePrintouts,
	}

	return &sqlSvr
}

func (s *MssqlServer) Initialise() error {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)

	if err := s.createDashboarDatadDbIfNotExist(); err != nil {
		return err
	}

	if err := s.createUserTableIfNotExist(); err != nil {
		return err
	}
	if err := s.createTaskTableIfNotExist(); err != nil {
		return err
	}

	return nil
}

func (s *MssqlServer) establishConnection() (*sql.DB, error) {
	connectionString := s.generateConnectionString()
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return db, nil
}

func (s *MssqlServer) generateConnectionString() string {
	connectionString := fmt.Sprintf("server=%s; port=%d; user id=%s; password=%s; database=%s; trustedConnection=%t; trustServerCertificate=%t", s.ServerName, s.Port, s.Username, s.Password, s.DatabaseName, s.TrustedConnection, s.TrustServerCertificate)
	return connectionString
}
