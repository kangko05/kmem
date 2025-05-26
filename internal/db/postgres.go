package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Postgres struct {
	conn *sql.DB
}

func Connect() (*Postgres, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		"localhost", 5432, "postgres", "rkdehddn12", "testdb", "disable")

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &Postgres{conn: conn}, nil
}

func (pg *Postgres) Ping() error {
	return pg.conn.Ping()
}

func (pg *Postgres) Close() error {
	return pg.conn.Close()
}
