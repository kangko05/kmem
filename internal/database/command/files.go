package command

import (
	"context"
	"fmt"
	"kmem/internal/database"
	"kmem/internal/models"
	"time"
)

func InsertFileMetadata(pg *database.Postgres, metadata models.FileMetadata) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := pg.Exec(ctx,
		`INSERT INTO files(filename,contenttype,storedpath,uploadedby,uploadedat,size,hash,archivepath) VALUES($1,$2,$3,$4,$5,$6,$7, $8)`,
		metadata.Filename, metadata.ContentType, metadata.StoredPath, metadata.UploadedBy, metadata.UploadedAt, metadata.Size, metadata.Hash, metadata.ArchivePath,
	)
	if err != nil {
		return fmt.Errorf("failed to save metadata: %v", err)
	}

	return nil
}
