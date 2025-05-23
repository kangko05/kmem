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

func GetRecentitems(pg *database.Postgres, username, sortBy string, limit int) ([]models.MetadataPart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var queryStr = "SELECT storedpath,filename,contenttype,uploadedat FROM files WHERE uploadedby=$1"

	switch sortBy {
	case "name":
		queryStr += " ORDER BY filename ASC"

	case "date":
		fallthrough
	default:
		queryStr += " ORDER BY uploadedat DESC"
	}

	var result []models.MetadataPart

	scan := func(rows *sql.Rows) error {
		for rows.Next() {
			var mp models.MetadataPart

			if err := rows.Scan(&mp.Storedpath, &mp.Filename, &mp.ContentType, &mp.UploadedAt); err != nil {
				log.Printf("error scanning rows: %v\n", err)
				continue
			}

			result = append(result, mp)
		}

		return nil
	}

	err := pg.Query(ctx, scan, queryStr, username)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}

	return result, nil
}

func QueryPatternMatching(pg *database.Postgres, username, searchQuery string) ([]models.FileMetadata, error) {
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
		"SELECT filename,contenttype,storedpath,archivepath,uploadedby,uploadedat,size FROM files WHERE uploadedby=$1 AND (filename ILIKE $2 OR contenttype ILIKE $2)",
		username, fmt.Sprintf("%%%s%%", searchQuery),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query user file: %v", err)
	}

	return files, nil
}

// get a local path to single item
func GetItem(pg *database.Postgres, username, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var result string

	scan := func(row *sql.Row) error {
		if err := row.Scan(&result); err != nil {
			return err
		}

		return nil
	}

	err := pg.QueryRow(ctx, scan, "SELECT storedpath FROM files WHERE uploadedby=$1 AND filename=$2", username, filename)
	if err != nil {
		return "", err
	}

	return result, nil
}
