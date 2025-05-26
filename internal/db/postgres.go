package db

import (
	"context"
	"database/sql"
	"fmt"
	"kmem/internal/config"
	"time"

	_ "github.com/lib/pq"
)

type Postgres struct {
	conn *sql.DB
}

func Connect(conf *config.Config) (*Postgres, error) {
	conn, err := sql.Open("postgres", conf.PostgresConnStr())
	if err != nil {
		return nil, err
	}

	pg := &Postgres{conn: conn}

	if err := pg.initTables(); err != nil {
		return nil, fmt.Errorf("failed to init tables: %v", err)
	}

	return pg, nil
}

// users
func (pg *Postgres) initTables() error {
	return pg.Exec(`CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		username VARCHAR(20) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_login TIMESTAMP
	)`)
}

func (pg *Postgres) Ping() error {
	return pg.conn.Ping()
}

func (pg *Postgres) Close() error {
	return pg.conn.Close()
}

func (pg *Postgres) Exec(query string, args ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pg.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(query, args...); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return rerr
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
