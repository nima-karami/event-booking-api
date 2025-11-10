package models

import (
	"time"

	"example.com/event-booking-api/db"
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
	return err
}

func (e *Event) Update() error {
	query := `
    UPDATE events
    SET title = $1, description = $2, location = $3, date = $4
    WHERE id = $5
    `
	_, err := db.DB.Exec(query, e.Title, e.Description, e.Location, e.Date, e.ID)
	return err
}

func (e *Event) Delete() error {
	query := "DELETE FROM events WHERE id = $1"
	_, err := db.DB.Exec(query, e.ID)
	return err
}

func (e *Event) Register(userID int64) error {
	query := "INSERT INTO registrations (user_id, event_id) VALUES ($1, $2)"
	_, err := db.DB.Exec(query, userID, e.ID)
	return err
}

func (e *Event) Unregister(userID int64) error {
	query := "DELETE FROM registrations WHERE user_id = $1 AND event_id = $2"
	_, err := db.DB.Exec(query, userID, e.ID)
	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT id, title, description, location, date, user_id FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.Date, &e.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := "SELECT id, title, description, location, date, user_id FROM events WHERE id = $1"
	row := db.DB.QueryRow(query, id)

	var e Event
	err := row.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.Date, &e.UserID)
	if err != nil {
		return nil, err
	}

	return &e, nil
}
