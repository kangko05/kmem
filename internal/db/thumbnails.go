package db

import (
	"fmt"
	"kmem/internal/models"
)

func (pg *Postgres) InsertThumbnails(t models.Thumbnail) error {
	err := pg.Exec(`
		INSERT INTO thumbnails(file_id,size_name,width,height,file_path,relative_path,file_size)
		VALUES($1,$2,$3,$4,$5,$6,$7)
		`, t.FileID, t.SizeName, t.Width, t.Height, t.FilePath, t.RelativePath, t.FileSize)
	if err != nil {
		return fmt.Errorf("failed to insert thumbnail: %v", err)
	}

	return nil
}
