package db

import (
	"context"
	"database/sql"
	"fmt"
	"kmem/internal/models"
	"log"
	"time"
)

func (pg *Postgres) InsertFile(file models.File) (int, error) {
	// tx
	txctx, cancel := context.WithTimeout(pg.ctx, pg.txtimeout)
	defer cancel()

	tx, err := pg.conn.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer tx.Rollback()

	// check existing file
	var id int
	var deleted bool
	err = tx.QueryRowContext(txctx, `SELECT id,deleted FROM files WHERE hash=$1`, file.Hash).Scan(&id, &deleted)
	if err == nil { // file exists
		fmt.Println(id, deleted)

		if deleted {
			if _, err := tx.ExecContext(txctx, `UPDATE files SET deleted=$1,deleted_at=$2 WHERE id=$3`, false, nil, id); err != nil {
				fmt.Println(err)
				return 0, fmt.Errorf("failed to update deleted file: %v", err)
			}

			if err := tx.Commit(); err != nil {
				return 0, fmt.Errorf("failed to commit tx: %v", err)
			}

			return id, nil
		} else {
			return id, fmt.Errorf("file already exists")
		}
	}

	// new file
	err = tx.QueryRowContext(txctx, `
	INSERT INTO files(username,hash,original_name,stored_name,file_path,relative_path,file_size,mime_type)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING id
	`, file.Username, file.Hash, file.OriginalName, file.StoredName, file.FilePath, file.RelativePath, file.FileSize, file.MimeType).Scan(&id)
	if err != nil {
		return 0, err
	}

	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	// 	return -1, err
	// }
	//
	// if rowsAffected == 0 {
	// 	return -1, fmt.Errorf("file already exists")
	// }

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit tx: %v", err)
	}

	return id, nil
}

func (pg *Postgres) GetFilesCount(username, typeStr, searchStr string) (int, error) {
	whereClause := "WHERE username=$1"
	args := []any{username}

	if typeStr != "all" {
		whereClause += " AND mime_type LIKE $2"
		args = append(args, typeStr+"%")
	}

	if len(searchStr) > 2 {
		whereClause += fmt.Sprintf(" AND original_name ILIKE $%d", len(args)+1)
		args = append(args, "%"+searchStr+"%")
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM files %s", whereClause)

	var count int
	err := pg.conn.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("failed to get file count: %v", err)
	}
	return count, nil
}

func (pg *Postgres) GetFilesPage(username string, page, limit int, sort, typeStr, searchStr string) ([]models.FileResponse, error) {
	offset := page * limit

	orderby := "uploaded_at DESC" // default
	if sort == "date" {
		orderby = "uploaded_at DESC"
	} else if sort == "name" {
		orderby = "original_name ASC"
	}

	whereClause := "WHERE username=$1 AND deleted=$2"
	args := []any{username, false}

	if typeStr != "all" {
		whereClause += " AND mime_type LIKE $3"
		args = append(args, typeStr+"%")
	}

	if len(searchStr) > 2 {
		whereClause += fmt.Sprintf(" AND original_name ILIKE $%d", len(args)+1)
		args = append(args, "%"+searchStr+"%")
	}

	query := fmt.Sprintf(`
		SELECT f.id,f.original_name,f.relative_path,f.mime_type,t.size_name,t.relative_path FROM (
			SELECT id, original_name, relative_path, mime_type
        	FROM files 
        	%s
        	ORDER BY %s
        	LIMIT $%d OFFSET $%d
		) AS f
		LEFT JOIN thumbnails AS t ON f.id=t.file_id
	`, whereClause, orderby, len(args)+1, len(args)+2)

	args = append(args, limit, offset)

	rows, err := pg.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get files for %s: %v", username, err)
	}
	defer rows.Close()

	filesMap := make(map[int]models.FileResponse)
	for rows.Next() {
		var file models.FileResponse

		file.Thumbnails = make(map[string]models.ThumbnailResponse)

		var sizeName, thumbPath sql.NullString

		if err := rows.Scan(&file.ID, &file.OriginalName, &file.FilePath, &file.MimeType, &sizeName, &thumbPath); err != nil {
			log.Println(err)
			continue
		}

		var thumb models.ThumbnailResponse
		if sizeName.Valid && thumbPath.Valid {
			thumb.SizeName = sizeName.String
			thumb.FilePath = thumbPath.String
		}

		_, ok := filesMap[file.ID]
		if ok {
			filesMap[file.ID].Thumbnails[thumb.SizeName] = thumb
		} else {
			filesMap[file.ID] = file
		}
	}

	var files []models.FileResponse

	for _, f := range filesMap {
		files = append(files, f)
	}

	return files, nil
}

// soft remove files - local files will be deleted after some time
func (pg *Postgres) DeleteFile(username, fileId string) error {
	txctx, cancel := context.WithTimeout(pg.ctx, pg.txtimeout)
	defer cancel()

	tx, err := pg.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin tx: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(txctx, `
	UPDATE files
	SET deleted=$1,deleted_at=$2
	WHERE username=$3 AND id=$4`,
		true, time.Now(), username, fileId)
	if err != nil {
		return fmt.Errorf("failed to delete soft file from db: %s: %v", fileId, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %v", err)
	}

	return nil
}

func (pg *Postgres) RenameFile(username, fileId, newName string) error {
	txctx, cancel := context.WithTimeout(pg.ctx, pg.txtimeout)
	defer cancel()

	tx, err := pg.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin tx: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(txctx, `UPDATE files SET original_name=$1 WHERE username=$2 AND id=$3 AND deleted=$4`, newName, username, fileId, false)
	if err != nil {
		return fmt.Errorf("failed to rename file: %s: %v", fileId, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %v", err)
	}

	return nil
}
