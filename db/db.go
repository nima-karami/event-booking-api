package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic(err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createEventsTableSQL := `CREATE TABLE IF NOT EXISTS events (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT NOT NULL,
		"description" TEXT NOT NULL,
		"location" TEXT NOT NULL,
		"date" DATETIME NOT NULL,
		"user_id" INTEGER
	);`
	_, err := DB.Exec(createEventsTableSQL)
	if err != nil {
		panic(err)
	}
}
