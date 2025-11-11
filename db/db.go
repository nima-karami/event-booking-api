package db

import (
	"database/sql"
	"fmt"

	"example.com/event-booking-api/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	utils.Logger.Debug("Connecting to database",
		"host", host,
		"port", port,
		"database", dbname)

	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		utils.Logger.Error("Failed to open database connection", "error", err)
		panic("Failed to connect to database: " + err.Error())
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	err = DB.Ping()
	if err != nil {
		utils.Logger.Error("Failed to ping database", "error", err)
		panic("Could not ping database: " + err.Error())
	}

	utils.Logger.Info("Database connection established",
		"host", host,
		"port", port,
		"database", dbname)

	runMigrations(connStr)
}

func runMigrations(connStr string) {
	utils.Logger.Info("Running database migrations")

	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		utils.Logger.Error("Failed to create migration driver", "error", err)
		panic("Could not create migration driver: " + err.Error())
	}

	migrationsPath := utils.GetEnvString("MIGRATIONS_PATH", "file://db/migrations")

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		utils.Logger.Error("Failed to initialize migrations", "error", err, "path", migrationsPath)
		panic("Could not initialize migrations: " + err.Error())
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		utils.Logger.Error("Failed to run migrations", "error", err)
		panic("Could not run migrations: " + err.Error())
	}

	if err == migrate.ErrNoChange {
		utils.Logger.Info("No new migrations to apply")
	} else {
		utils.Logger.Info("Migrations applied successfully")
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		utils.Logger.Warn("Failed to get migration version", "error", err)
	} else if err == nil {
		utils.Logger.Info("Current migration version", "version", version, "dirty", dirty)
	}
}
