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
	createUsersTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"email" TEXT NOT NULL UNIQUE,
		"password" TEXT NOT NULL
	);`
	_, err := DB.Exec(createUsersTableSQL)
	if err != nil {
		panic(err)
	}

	createEventsTableSQL := `CREATE TABLE IF NOT EXISTS events (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT NOT NULL,
		"description" TEXT NOT NULL,
		"location" TEXT NOT NULL,
		"date" DATETIME NOT NULL,
		"user_id" INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`
	_, err = DB.Exec(createEventsTableSQL)
	if err != nil {
		panic(err)
	}

	createRegistrationsTableSQL := `CREATE TABLE IF NOT EXISTS registrations (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"user_id" INTEGER,
		"event_id" INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(event_id) REFERENCES events(id)
	);`
	_, err = DB.Exec(createRegistrationsTableSQL)
	if err != nil {
		panic(err)
	}
}
