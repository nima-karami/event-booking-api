package db

import (
	"database/sql"
	"fmt"

	"example.com/event-booking-api/utils"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	host := utils.GetEnvString("DB_HOST", "localhost")
	port := utils.GetEnvString("DB_PORT", "5432")
	user := utils.GetEnvString("DB_USER", "postgres")
	password := utils.GetEnvString("DB_PASSWORD", "postgres")
	dbname := utils.GetEnvString("DB_NAME", "event_booking")
	sslmode := utils.GetEnvString("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	fmt.Printf("Connecting to database with connection string: %s\n", connStr)

	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );
    `

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		panic("Could not create users table: " + err.Error())
	}

	createEventsTable := `
    CREATE TABLE IF NOT EXISTS events (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        description TEXT NOT NULL,
        location TEXT NOT NULL,
        date TIMESTAMP NOT NULL,
        user_id INTEGER NOT NULL,
        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
    );
    `

	_, err = DB.Exec(createEventsTable)
	if err != nil {
		panic("Could not create events table: " + err.Error())
	}

	createRegistrationsTable := `
    CREATE TABLE IF NOT EXISTS registrations (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        event_id INTEGER NOT NULL,
        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE,
        UNIQUE(user_id, event_id)
    );
    `

	_, err = DB.Exec(createRegistrationsTable)
	if err != nil {
		panic("Could not create registrations table: " + err.Error())
	}
}
