package db

import (
	"database/sql"
	"log"

	// import pq driver
	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDatabase connect DB
func InitDatabase(databaseURL string) *sql.DB {
	var err error
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

// GetConnection connect DB
func GetConnection() *sql.DB {
	return db
}
