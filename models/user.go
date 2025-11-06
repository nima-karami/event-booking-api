package models

import (
	"example.com/event-booking-api/db"
	"example.com/event-booking-api/utils"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) Save() error {
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (email, password) VALUES (?, ?)`
	result, err := db.DB.Exec(query, u.Email, hashedPassword)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	u.ID = id
	return err
}

func (u *User) Update() error {
	query := `UPDATE users SET email = ?, password = ? WHERE id = ?`
	_, err := db.DB.Exec(query, u.Email, u.Password, u.ID)
	return err
}

func (u *User) Delete() error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := db.DB.Exec(query, u.ID)
	return err
}

func (u *User) Authenticate() error {
	query := `SELECT id, password FROM users WHERE email = ?`
	row := db.DB.QueryRow(query, u.Email)

	var storedHashedPassword string
	err := row.Scan(&u.ID, &storedHashedPassword)
	if err != nil {
		return err
	}

	return utils.CheckPasswordHash(u.Password, storedHashedPassword)
}

func GetAllUsers() ([]User, error) {
	query := `SELECT * FROM users`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Email, &u.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func GetUserByID(id int64) (*User, error) {
	query := `SELECT * FROM users WHERE id = ?`
	row := db.DB.QueryRow(query, id)

	var u User
	err := row.Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
