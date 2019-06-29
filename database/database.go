package database

import (
	"database/sql"
	"os"
)

type Interface interface {
	Connect() (*sql.DB, error)
}

type PostgresDB struct {
}

func (p PostgresDB) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return db, nil
}
