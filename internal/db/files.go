package db

import (
	"database/sql"
	"fmt"
	"kmem/internal/models"
	"log"
)

func (pg *Postgres) InsertFile(file models.File) (int, error) {
	var fileId int
	err := pg.conn.QueryRow(`
	INSERT INTO files(username,hash,original_name,stored_name,file_path,relative_path,file_size,mime_type)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING id
	`, file.Username, file.Hash, file.OriginalName, file.StoredName, file.FilePath, file.RelativePath, file.FileSize, file.MimeType).Scan(&fileId)

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

	return fileId, nil
}

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

func (pg *Postgres) GetFilesPage(username string, page, limit int, sort, typeStr string) ([]models.FileResponse, error) {
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
		var id int
		var file models.FileResponse

		file.Thumbnails = make(map[string]models.ThumbnailResponse)

		var sizeName, thumbPath sql.NullString

		if err := rows.Scan(&id, &file.OriginalName, &file.FilePath, &file.MimeType, &sizeName, &thumbPath); err != nil {
			log.Println(err)
			continue
		}

		var thumb models.ThumbnailResponse
		if sizeName.Valid && thumbPath.Valid {
			thumb.SizeName = sizeName.String
			thumb.FilePath = thumbPath.String
		}

		_, ok := filesMap[id]
		if ok {
			filesMap[id].Thumbnails[thumb.SizeName] = thumb
		} else {
			filesMap[id] = file
		}
	}

	var files []models.FileResponse

	for _, f := range filesMap {
		files = append(files, f)
	}

	return files, nil
}

// func (pg *Postgres) GetFiles(username string) ([]models.File, error) {
// 	rows, err := pg.conn.Query(`
// 		SELECT original_name,file_path,mime_type FROM files
// 		WHERE username=$1
// 	`, username)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get files for %s: %v", username, err)
// 	}
// 	defer rows.Close()
//
// 	var files []models.File
// 	for rows.Next() {
// 		var file models.File
// 		if err := rows.Scan(&file.OriginalName, &file.FilePath, &file.MimeType); err != nil {
// 			log.Println(err)
// 			continue
// 		}
//
// 		files = append(files, file)
// 	}
//
// 	return files, nil
// }
