package db

import (
	"fmt"
	"kmem/internal/models"
	"log"
)

func (pg *Postgres) InsertFile(file models.File) error {
	result, err := pg.conn.Exec(`
	INSERT INTO files(username,hash,original_name,stored_name,file_path,file_size,mime_type)
	VALUES($1,$2,$3,$4,$5,$6,$7)
	ON CONFLICT (hash) DO NOTHING
	`, file.Username, file.Hash, file.OriginalName, file.StoredName, file.FilePath, file.FileSize, file.MimeType)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("file already exists")
	}

	return nil
}

// type File struct {
// 	ID           int       `json:"id" db:"id"`
// 	Hash         string    `json:"hash" db:"hash"` // md5
// 	Username     string    `json:"username" db:"username"`
// 	OriginalName string    `json:"originalName" db:"original_name"`
// 	StoredName   string    `json:"storedName" db:"stored_name"`
// 	FilePath     string    `json:"filePath" db:"file_path"`
// 	FileSize     int64     `json:"fileSize" db:"file_size"`
// 	MimeType     string    `json:"mimeType" db:"mime_type"`
// 	UploadedAt   time.Time `json:"uploadedAt" db:"uploaded_at"`
// }

func (pg *Postgres) GetFilesCount(username, typeStr string) (int, error) {
	whereClause := "WHERE username=$1"
	args := []any{username}

	if typeStr != "all" {
		whereClause += " AND mime_type LIKE $2"
		args = append(args, typeStr+"%")
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM files %s", whereClause)

	var count int
	err := pg.conn.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("failed to get file count: %v", err)
	}
	return count, nil
}

func (pg *Postgres) GetFilesPage(username string, page, limit int, sort, typeStr string) ([]models.File, error) {
	offset := page * limit

	orderby := "uploaded_at DESC" // default
	if sort == "date" {
		orderby = "uploaded_at DESC"
	} else if sort == "name" {
		orderby = "original_name ASC"
	}

	whereClause := "WHERE username=$1"
	args := []any{username}

	if typeStr != "all" {
		whereClause += " AND mime_type LIKE $2"
		args = append(args, typeStr+"%")
	}

	query := fmt.Sprintf(`
		SELECT original_name,file_path,mime_type FROM files
		%s
		ORDER BY %s
		LIMIT $%d
		OFFSET $%d
	`, whereClause, orderby, len(args)+1, len(args)+2)

	args = append(args, limit, offset)

	rows, err := pg.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get files for %s: %v", username, err)
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.OriginalName, &file.FilePath, &file.MimeType); err != nil {
			log.Println(err)
			continue
		}
		files = append(files, file)
	}

	return files, nil
}

func (pg *Postgres) GetFiles(username string) ([]models.File, error) {
	rows, err := pg.conn.Query(`
		SELECT original_name,file_path,mime_type FROM files
		WHERE username=$1
	`, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get files for %s: %v", username, err)
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.OriginalName, &file.FilePath, &file.MimeType); err != nil {
			log.Println(err)
			continue
		}

		files = append(files, file)
	}

	return files, nil
}
