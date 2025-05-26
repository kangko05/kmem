package db

import (
	"fmt"
	"kmem/internal/models"
	"kmem/internal/utils"
)

func (pg *Postgres) InsertUser(user models.User) error {
	_, err := pg.QueryUser(user.Username)
	if err == nil {
		return fmt.Errorf("username already exists")
	}

	hashedPass, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash pasword: %v", err)
	}

	return pg.Exec(`INSERT INTO users(username,password) VALUES($1,$2)`, user.Username, hashedPass)
}

func (pg *Postgres) QueryUser(username string) (models.User, error) {
	var user models.User

	err := pg.conn.QueryRow(`SELECT username,password FROM users WHERE username=$1`, username).Scan(&user.Username, &user.Password)
	if err != nil {
		return user, fmt.Errorf("failed to query user: %v", err)
	}

	return user, nil
}
