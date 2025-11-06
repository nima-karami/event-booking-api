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
	UserID      int       `json:"user_id"`
}

func (e *Event) Save() error {
	query := `INSERT INTO events (title, description, location, date, user_id) VALUES (?, ?, ?, ?, ?)`

	result, err := db.DB.Exec(query, e.Title, e.Description, e.Location, e.Date, e.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	e.ID = id
	return err

}

func (e *Event) Update() error {
	query := `UPDATE events SET title = ?, description = ?, location = ?, date = ? WHERE id = ?`
	_, err := db.DB.Exec(query, e.Title, e.Description, e.Location, e.Date, e.ID)
	return err
}

func (e *Event) Delete() error {
	query := `DELETE FROM events WHERE id = ?`
	_, err := db.DB.Exec(query, e.ID)
	return err
}

func GetAllEvents() ([]Event, error) {
	query := `SELECT * FROM events`
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
	query := `SELECT * FROM events WHERE id = ?`
	row := db.DB.QueryRow(query, id)

	var e Event
	err := row.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.Date, &e.UserID)
	if err != nil {
		return nil, err
	}

	return &e, nil
}
