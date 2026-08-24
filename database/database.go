package database

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

var database *sql.DB = nil

func GetDatabase() *sql.DB {
	if database != nil {
		return database
	}
	database, err := sql.Open("sqlite", "database/trading-card-album.sqlite3")
	if err != nil {
		panic("[DATABASE CONNECTION ERROR ] " + err.Error())
	}
	return database
}
