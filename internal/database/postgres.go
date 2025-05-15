package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Postgres struct {
	conn *sql.DB
}

func Connect(ctx context.Context) (*Postgres, error) {
	connStr := fmt.Sprintf("host=localhost port=5432 user=kang password=%s dbname=kmem sslmode=disable", "rkdehddn12")
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &Postgres{conn: conn}, nil
}

func (pg *Postgres) Close() {
	pg.conn.Close()
}

// for commands
func (pg *Postgres) Exec(ctx context.Context, query string, args ...any) error {
	tx, err := pg.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %v", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback: %v", rbErr)
		}

		return fmt.Errorf("failed to exec query: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %v", err)
	}

	return nil
}
