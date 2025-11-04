package db

import (
	"database/sql"

	"github.com/rs/zerolog/log"
)

var DB_DRIVER string

func init() {
	
}

func DBClient(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open(DB_DRIVER, dbUrl)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Print("Connected to DB!")
	return db, nil
}
