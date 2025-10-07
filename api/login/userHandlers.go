package login

import (
	"bitresume/config"
	"database/sql"
	"errors"
)

type User struct {
	Email  string
	RollNo string
	Role   string
	UserName string
}

// GetUserByEmail checks if a user exists and returns their details
func GetUserByEmail(email string) (*User, error) {
	var user User
	query := `SELECT user_email, rollno, role,user_name FROM login WHERE user_email = ?`
	err := config.DB.QueryRow(query, email).Scan(&user.Email, &user.RollNo, &user.Role ,&user.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
