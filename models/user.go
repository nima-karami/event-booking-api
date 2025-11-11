package models

import (
	"example.com/event-booking-api/db"
	"example.com/event-booking-api/utils"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

type PublicUser struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (u *User) Save() error {
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		utils.Logger.Error("Failed to hash password", "email", u.Email, "error", err)
		return err
	}

	query := `INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id`
	u.Role = "user" // Default role
	err = db.DB.QueryRow(query, u.Email, hashedPassword, u.Role).Scan(&u.ID)
	if err != nil {
		utils.Logger.Error("Failed to save user to database", "email", u.Email, "error", err)
		return err
	}
	utils.Logger.Debug("User saved to database", "user_id", u.ID, "email", u.Email)
	return nil
}

func (u *User) Update() error {
	query := `UPDATE users SET email = $1, password = $2 WHERE id = $3`
	_, err := db.DB.Exec(query, u.Email, u.Password, u.ID)
	if err != nil {
		utils.Logger.Error("Failed to update user in database", "user_id", u.ID, "email", u.Email, "error", err)
		return err
	}
	utils.Logger.Debug("User updated in database", "user_id", u.ID, "email", u.Email)
	return nil
}

func (u *User) Delete() error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.DB.Exec(query, u.ID)
	if err != nil {
		utils.Logger.Error("Failed to delete user from database", "user_id", u.ID, "error", err)
		return err
	}
	utils.Logger.Debug("User deleted from database", "user_id", u.ID)
	return nil
}

func (u *User) Authenticate() error {
	query := "SELECT id, password, role FROM users WHERE email = $1"
	row := db.DB.QueryRow(query, u.Email)

	var storedHashedPassword string
	err := row.Scan(&u.ID, &storedHashedPassword, &u.Role)
	if err != nil {
		utils.Logger.Error("Failed to retrieve user for authentication", "email", u.Email, "error", err)
		return err
	}

	err = utils.CheckPasswordHash(u.Password, storedHashedPassword)
	if err != nil {
		utils.Logger.Warn("Password verification failed", "email", u.Email)
		return err
	}

	utils.Logger.Debug("User authenticated successfully", "user_id", u.ID, "email", u.Email)
	return nil
}

func (u *User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:    u.ID,
		Email: u.Email,
		Role:  u.Role,
	}
}

func GetAllUsers() ([]User, error) {
	query := `SELECT * FROM users`
	rows, err := db.DB.Query(query)
	if err != nil {
		utils.Logger.Error("Failed to query all users", "error", err)
		return nil, err
	}

	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Role)
		if err != nil {
			utils.Logger.Error("Failed to scan user row", "error", err)
			return nil, err
		}
		users = append(users, u)
	}

	utils.Logger.Debug("Retrieved all users from database", "count", len(users))
	return users, nil
}

func GetUserByID(id int64) (*User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	row := db.DB.QueryRow(query, id)

	var u User
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role)
	if err != nil {
		utils.Logger.Error("Failed to get user by ID", "user_id", id, "error", err)
		return nil, err
	}

	utils.Logger.Debug("Retrieved user by ID", "user_id", id, "email", u.Email)
	return &u, nil
}

func GetUserByEmail(email string) (*User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	row := db.DB.QueryRow(query, email)

	var u User
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role)
	if err != nil {
		utils.Logger.Debug("User not found by email", "email", email, "error", err)
		return nil, err
	}

	utils.Logger.Debug("Retrieved user by email", "user_id", u.ID, "email", email)
	return &u, nil
}
