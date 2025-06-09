package queue

import (
	"fmt"
	"io/fs"
	"kmem/internal/cache"
	"kmem/internal/config"
	"kmem/internal/db"
	"kmem/internal/models"
	"log"
	"os"
	"path/filepath"
	"time"
)

type cleanItems struct {
	pg          *db.Postgres
	conf        *config.Config
	cache       *cache.Cache
	deleteAfter time.Duration
}

func CleanItems(pg *db.Postgres, conf *config.Config, cache *cache.Cache) *cleanItems {
	c := &cleanItems{
		pg:          pg,
		conf:        conf,
		cache:       cache,
		deleteAfter: time.Hour * 24 * 7 * 30,
	}

	return c
}

func (c *cleanItems) handleDeletedFile(dfile models.DelFile) error {
	now := time.Now()
	if !c.shouldDelete(now, *dfile.DeletedAt) {
		return nil
	}

	// remove from db first
	if err := c.pg.DeleteFileHard(dfile.Id); err != nil {
		return fmt.Errorf("failed to delete file hard: %d: %v", dfile.Id, err)
	}

	// remove all local files
	if err := os.Remove(dfile.FilePath); err != nil {
		log.Printf("failed to remove file %s: %v", dfile.FilePath, err)
	}

	for _, thumb := range dfile.ThumbnailPaths {
		if err := os.Remove(thumb); err != nil {
			log.Printf("failed to remove thumbnail %s: %v", thumb, err)
		}
	}

	return nil
}

func (c *cleanItems) checkFile(dfile models.DelFile) error {
	_, err := os.Stat(dfile.FilePath)
	if err != nil && os.IsNotExist(err) {
		if derr := c.pg.DeleteFileHard(dfile.Id); derr != nil {
			return fmt.Errorf("failed to delete orphaned file data %d: %v", dfile.Id, derr)
		}
	}
	return nil
}

func (c *cleanItems) checkLocalFiles(dmap map[string]models.DelFile) error {
	dbPaths := make(map[string]bool)
	for filePath, dfile := range dmap {
		dbPaths[filePath] = true
		for _, thumbPath := range dfile.ThumbnailPaths {
			if thumbPath != "" {
				dbPaths[thumbPath] = true
			}
		}
	}

	return filepath.WalkDir(c.conf.UploadPath(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if !dbPaths[path] {
			log.Printf("Removing orphaned file: %s", path)
			if err := os.Remove(path); err != nil {
				log.Printf("Failed to remove orphaned file %s: %v", path, err)
			}
		}

		return nil
	})
}

func (c *cleanItems) process() error {
	dmap, err := c.pg.GetAllFilesToCheck()
	if err != nil {
		return fmt.Errorf("failed to get files to clean: %v", err)
	}

	if err := c.checkLocalFiles(dmap); err != nil {
		log.Printf("something wrong while checking local files: %v\n", err)
	}

	for _, dfile := range dmap {
		if dfile.Deleted {
			err := c.handleDeletedFile(dfile)
			if err != nil {
				log.Println(err)
				continue
			}
		} else {
			err := c.checkFile(dfile)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}

	// clear gallery cache
	c.cache.ClearGalleryCache()

	return nil
}

func (c *cleanItems) shouldDelete(now, deletedAt time.Time) bool {
	return deletedAt.Add(c.deleteAfter).Unix() < now.Unix()
}
