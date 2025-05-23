package query

import (
	"context"
	"database/sql"
	"fmt"
	"kmem/internal/database"
	"kmem/internal/models"
	"log"
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

func QueryUsernames(pg *database.Postgres) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var usernames []string
	scan := func(rows *sql.Rows) error {
		var username string
		for rows.Next() {
			if err := rows.Scan(&username); err != nil {
				log.Printf("failed to scan username: %v", err)
				continue
			}

			usernames = append(usernames, username)
		}

		return nil
	}

	if err := pg.Query(ctx, scan, "SELECT username FROM users"); err != nil {
		return nil, fmt.Errorf("failed to query usernames: %v", err)
	}

	return usernames, nil
}

func QueryUserFiles(pg *database.Postgres, username string) ([]models.FileMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var files []models.FileMetadata
	scan := func(rows *sql.Rows) error {
		var file models.FileMetadata

		for rows.Next() {
			if err := rows.Scan(
				&file.Filename,
				&file.ContentType,
				&file.StoredPath,
				&file.ArchivePath,
				&file.UploadedBy,
				&file.UploadedAt,
				&file.Size,
			); err != nil {
				log.Printf("error scanning row: %v", err)
				continue
			}

			files = append(files, file)
		}

		return nil
	}

	err := pg.Query(
		ctx,
		scan,
		"SELECT filename,contenttype,storedpath,archivepath,uploadedby,uploadedat,size FROM files LEFT JOIN users ON files.uploadedby=users.username WHERE users.username=$1",
		username,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query user file: %v", err)
	}

	return files, nil
}
