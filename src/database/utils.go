package database

import (
	"database/sql"
	"fmt"
	"log"
)

/*
	Exec SQL Commands
*/
func (s *MssqlServer) execQuery(query string) (*sql.Rows, error) {
	db, err := s.establishConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	return rows, nil
}

func (s *MssqlServer) execCommand(command string) error {
	db, err := s.establishConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(command)
	if err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}
	return nil
}

func (s *MssqlServer) execNamedCommand(namedCommand string, namedArgs ...interface{}) error {
	db, err := s.establishConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(namedCommand)
	if err != nil {
		return fmt.Errorf("error preparing named command: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedArgs...)
	if err != nil {
		return fmt.Errorf("error executing named command: %w", err)
	}
	return nil
}

func (s *MssqlServer) execNamedQuery(namedQuery string, namedArgs ...interface{}) (*sql.Rows, error) {
	db, err := s.establishConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(namedQuery)
	if err != nil {
		return nil, fmt.Errorf("error preparing named query: %w", err)
	}
	defer stmt.Close()
	
	rows, err := stmt.Query(namedArgs...)
	if err != nil {
		return nil, fmt.Errorf("error executing named query: %w", err)
	}
	return rows, nil
}

func (s *MssqlServer) printRows(rows *sql.Rows) {
	// Print Headers
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(columns)

	// Prepare a slice of interfaces to hold the values for each column
	values := make([]interface{}, len(columns))
	valuePointers := make([]interface{}, len(columns))

	for rows.Next() {
		for i := range values {
			valuePointers[i] = &values[i]
		}

		if err = rows.Scan(valuePointers...); err != nil {
			log.Fatal(fmt.Errorf("error scanning for rows: %w", err))
		}

		for _, val := range values {
			var v string
			// Check for byte array - Strings
			// Else, any other integral type
			if b, ok := val.([]byte); ok {
				v = string(b)
			} else {
				v = fmt.Sprintf("%v", val)
			}
			fmt.Print(v, ", ")
		}
		fmt.Println()
	}
	fmt.Println()
}