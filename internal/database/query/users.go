package query

import (
	"context"
	"database/sql"
	"fmt"
	"kmem/internal/database"
	"kmem/internal/models"
	"time"
)

func QueryUser(pg *database.Postgres, username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	scan := func(row *sql.Row) error {
		if err := row.Scan(&user.Username, &user.Password); err != nil {
			return err
		}

		return nil
	}

	err := pg.QueryRow(ctx, scan, "SELECT username,password FROM users WHERE username=$1", username)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to find user: %v", err)
	}

	return user, nil
}
