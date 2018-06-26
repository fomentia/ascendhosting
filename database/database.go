package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"

	"github.com/fomentia/ascendhosting/models"
)

var schema = `CREATE TABLE IF NOT EXISTS hosts (
  id SERIAL PRIMARY KEY,
  first_name VARCHAR(255) NOT NULL,
  last_name  VARCHAR(255) NOT NULL
)`

type DB struct {
	database *sql.DB
}

func InitDB() (*DB, error) {
	datastoreName := os.Getenv("ASCEND_HOSTING_POSTGRES_CONNECTION")

	database, err := sql.Open("postgres", datastoreName)
	if err != nil {
		return nil, err
	}

	_, err = database.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &DB{database}, nil
}

func (db *DB) Insert(model models.Model) (errors models.Errors) {
	errors = models.Validate(model)

	if errors.None() {
		err := model.Insert(db.database)

		// TODO: separate database errors from validation errors.
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	return errors
}
