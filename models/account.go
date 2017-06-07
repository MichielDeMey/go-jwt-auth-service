package models

import (
	"log"
)

type LoginCredentials struct {
	Email    string
	Password string
}

type User struct {
	ID       int
	Email    string
	Password string
}

func GetUserByEmail(email string) (*User, error) {
	user := new(User)

	stmt, err := db.Prepare("SELECT id, email, password FROM users WHERE email = $1")
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.QueryRow(email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		log.Fatal(err)
	}

	return user, err
}

func CreateUser(credentials LoginCredentials) (*User, error) {
	insertedUser := new(User)

	err := db.QueryRow("INSERT INTO users(email, password) VALUES($1, crypt($2, gen_salt('bf', 8))) RETURNING id, email, password", credentials.Email, credentials.Password).Scan(&insertedUser.ID, &insertedUser.Email, &insertedUser.Password)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser, err
}
