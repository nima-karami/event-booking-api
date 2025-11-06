package models

import "example.com/event-booking-api/db"

type RegistrationWithUser struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	EventID int64  `json:"event_id"`
	Email   string `json:"email"`
}

func GetRegistrationsByEventIDWithUsers(eventID int64) ([]RegistrationWithUser, error) {
	query := `
        SELECT r.id, r.user_id, r.event_id, u.email 
        FROM registrations r
        JOIN users u ON r.user_id = u.id
        WHERE r.event_id = ?
    `
	rows, err := db.DB.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	registrations := []RegistrationWithUser{}
	for rows.Next() {
		var r RegistrationWithUser
		err := rows.Scan(&r.ID, &r.UserID, &r.EventID, &r.Email)
		if err != nil {
			return nil, err
		}
		registrations = append(registrations, r)
	}

	return registrations, nil
}
