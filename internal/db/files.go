package db

import (
	"fmt"
	"kmem/internal/models"
)

// func (pg *Postgres) InsertFile(file models.File) error {
// 	return pg.Exec(`
// 	INSERT INTO files(username,hash,original_name,stored_name,file_path,file_size,mime_type)
// 	VALUES($1,$2,$3,$4,$5,$6,$7)
// 	ON CONFLICT DO NOTHING
// 	`, file.Username, file.Hash, file.OriginalName, file.StoredName, file.FilePath, file.FileSize, file.MimeType)
// }

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
