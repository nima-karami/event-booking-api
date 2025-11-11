package models

import (
	"time"

	"example.com/event-booking-api/db"
	"example.com/event-booking-api/utils"
)

type Event struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	UserID      int64     `json:"user_id"`
}

func (e *Event) Save() error {
	query := `
    INSERT INTO events (title, description, location, date, user_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id
    `
	err := db.DB.QueryRow(query, e.Title, e.Description, e.Location, e.Date, e.UserID).Scan(&e.ID)
	if err != nil {
		utils.Logger.Error("Failed to save event to database",
			"title", e.Title,
			"user_id", e.UserID,
			"error", err)
		return err
	}
	utils.Logger.Debug("Event saved to database", "event_id", e.ID, "title", e.Title)
	return nil
}

func (e *Event) Update() error {
	query := `
    UPDATE events
    SET title = $1, description = $2, location = $3, date = $4
    WHERE id = $5
    `
	_, err := db.DB.Exec(query, e.Title, e.Description, e.Location, e.Date, e.ID)
	if err != nil {
		utils.Logger.Error("Failed to update event in database",
			"event_id", e.ID,
			"title", e.Title,
			"error", err)
		return err
	}
	utils.Logger.Debug("Event updated in database", "event_id", e.ID, "title", e.Title)
	return nil
}

func (e *Event) Delete() error {
	query := "DELETE FROM events WHERE id = $1"
	_, err := db.DB.Exec(query, e.ID)
	if err != nil {
		utils.Logger.Error("Failed to delete event from database",
			"event_id", e.ID,
			"error", err)
		return err
	}
	utils.Logger.Debug("Event deleted from database", "event_id", e.ID)
	return nil
}

func (e *Event) Register(userID int64) error {
	query := "INSERT INTO registrations (user_id, event_id) VALUES ($1, $2)"
	_, err := db.DB.Exec(query, userID, e.ID)
	if err != nil {
		utils.Logger.Error("Failed to register user for event",
			"event_id", e.ID,
			"user_id", userID,
			"error", err)
		return err
	}
	utils.Logger.Debug("User registered for event", "event_id", e.ID, "user_id", userID)
	return nil
}

func (e *Event) Unregister(userID int64) error {
	query := "DELETE FROM registrations WHERE user_id = $1 AND event_id = $2"
	_, err := db.DB.Exec(query, userID, e.ID)
	if err != nil {
		utils.Logger.Error("Failed to unregister user from event",
			"event_id", e.ID,
			"user_id", userID,
			"error", err)
		return err
	}
	utils.Logger.Debug("User unregistered from event", "event_id", e.ID, "user_id", userID)
	return nil
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT id, title, description, location, date, user_id FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		utils.Logger.Error("Failed to query all events", "error", err)
		return nil, err
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.Date, &e.UserID)
		if err != nil {
			utils.Logger.Error("Failed to scan event row", "error", err)
			return nil, err
		}
		events = append(events, e)
	}

	utils.Logger.Debug("Retrieved all events from database", "count", len(events))
	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := "SELECT id, title, description, location, date, user_id FROM events WHERE id = $1"
	row := db.DB.QueryRow(query, id)

	var e Event
	err := row.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.Date, &e.UserID)
	if err != nil {
		utils.Logger.Error("Failed to get event by ID", "event_id", id, "error", err)
		return nil, err
	}

	utils.Logger.Debug("Retrieved event by ID", "event_id", id, "title", e.Title)
	return &e, nil
}
