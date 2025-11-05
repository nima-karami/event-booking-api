package models

import "time"

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	UserID      int       `json:"user_id"`
}

var events = []Event{}

func (e *Event) Save() error {
	// Placeholder for saving the event to the database
	events = append(events, *e)
	return nil
}

func GetAllEvents() []Event {
	return events
}
