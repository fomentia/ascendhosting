package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var hostSchema = `CREATE TABLE IF NOT EXISTS hosts (
  first_name VARCHAR(255),
  last_name VARCHAR(255)
)`

type Host struct {
	firstName string
	lastName  string
}

type DB struct {
	database *sql.DB
}

func initDB() (*DB, error) {
	datastoreName := os.Getenv("ASCEND_HOSTING_POSTGRES_CONNECTION")

	database, err := sql.Open("postgres", datastoreName)
	if err != nil {
		return nil, err
	}

	_, err = database.Exec(hostSchema)
	if err != nil {
		return nil, err
	}

	return &DB{database}, nil
}

func (db *DB) insertHost(host Host) error {
	stmt := `INSERT INTO hosts (first_name, last_name) VALUES ($1, $2)`
	_, err := db.database.Exec(stmt, host.firstName, host.lastName)
	return err
}
