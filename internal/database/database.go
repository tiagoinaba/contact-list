package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func New() *sql.DB {
	db, err := sql.Open("sqlite3", "file:db.db?_foreign_keys=on")
	if err != nil {
		panic(err)
	}
	return db
}
