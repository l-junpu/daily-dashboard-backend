package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type MssqlServer struct {
	ServerName             string
	DatabaseName           string
	TrustedConnection      bool
	TrustServerCertificate bool
	EnablePrintouts        bool
}

/*
Connection Functions
*/
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
	connectionString := fmt.Sprintf("server=%s; database=%s; trustedConnection=%t; trustServerCertificate=%t", s.ServerName, s.DatabaseName, s.TrustedConnection, s.TrustServerCertificate)
	return connectionString
}
