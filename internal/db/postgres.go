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

func (pg *Postgres) initTables() error {
	err := pg.Exec(`CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		username VARCHAR(20) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_login TIMESTAMP
	)`)
	if err != nil {
		return fmt.Errorf("failed to init users table: %v", err)
	}

	err = pg.Exec(`CREATE TABLE IF NOT EXISTS files(
		id SERIAL PRIMARY KEY,
		hash VARCHAR(255) NOT NULL UNIQUE,
		username VARCHAR(20) NOT NULL,
		original_name VARCHAR(255) NOT NULL,
		stored_name VARCHAR(255) NOT NULL,
		file_path VARCHAR(255) NOT NULL,
		relative_path VARCHAR(255) NOT NULL,
		file_size BIGINT,
		mime_type VARCHAR(32),
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE
	)`)
	if err != nil {
		return fmt.Errorf("failed to init files table: %v", err)
	}

	err = pg.Exec(`CREATE TABLE IF NOT EXISTS thumbnails(
		id SERIAL PRIMARY KEY,
		file_id INTEGER NOT NULL,
		size_name VARCHAR(10) NOT NULL,
		width INTEGER NOT NULL,
		height INTEGER NOT NULL,
		file_path VARCHAR(255) NOT NULL,
		relative_path VARCHAR(255) NOT NULL,
		file_size BIGINT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(file_id,size_name),
		FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
	)`)
	if err != nil {
		return fmt.Errorf("failed to init thumbnails table: %v", err)
	}

	// TODO: add index

	return nil
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
