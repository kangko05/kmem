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
func (pg *Postgres) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	tx, err := pg.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("failed to rollback: %v", rbErr)
		}

		return nil, fmt.Errorf("failed to exec query: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %v", err)
	}

	return result, nil
}

// for queries
func (pg *Postgres) QueryRow(ctx context.Context, scanFn func(*sql.Row) error, query string, args ...any) error {
	tx, err := pg.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %v", err)
	}

	row := tx.QueryRowContext(ctx, query, args...)

	if err := scanFn(row); err != nil {
		return fmt.Errorf("failed to scan values: %v", err)
	}

	if err := tx.Commit(); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback: %v", rbErr)
		}

		return fmt.Errorf("failed to query row: %v", err)
	}

	return nil
}
