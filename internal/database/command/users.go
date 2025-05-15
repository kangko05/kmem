package command

import (
	"context"
	"fmt"
	"kmem/internal/database"
	"kmem/internal/models"
	"time"
)

// commands can mutate system state

func InsertUser(pg *database.Postgres, user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := pg.Exec(ctx, "INSERT INTO users(username, password) VALUES($1,$2)", user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	return nil
}

func UpdateLastLogin(pg *database.Postgres, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := pg.Exec(ctx, "UPDATE users SET last_login=$1 WHERE username=$2", time.Now(), username)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}

	// this shouldn't really happen
	if affected == 0 {
		return fmt.Errorf("username %s not found", username)
	}

	// this shouldn't really happen 2
	if affected > 1 {
		return fmt.Errorf("multiple rows found with username %s", username)
	}

	return nil
}
