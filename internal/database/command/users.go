package command

import (
	"context"
	"fmt"
	"kmem/internal/database"
	"kmem/internal/models"
	"time"
)

// commands can mutate system state

func StoreUser(pg *database.Postgres, user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := pg.Exec(ctx, "INSERT INTO users(username, password) VALUES($1,$2)", user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	return nil
}
