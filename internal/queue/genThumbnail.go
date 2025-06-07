package queue

import (
	"fmt"
	"kmem/internal/config"
	"kmem/internal/db"
	"kmem/internal/models"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type thumbnailSize struct {
	name   string // small, medium, large
	width  int
	height int
}

type genThumbnail struct {
	ts   []thumbnailSize
	file models.File
	pg   *db.Postgres
	conf *config.Config
}

func GenThumbnail(pg *db.Postgres, conf *config.Config, file models.File) *genThumbnail {
	return &genThumbnail{
		ts: []thumbnailSize{
			{name: "small", width: 150, height: 150},
			{name: "medium", width: 300, height: 300},
			{name: "large", width: 800, height: 600},
		},
		file: file,
		pg:   pg,
		conf: conf,
	}
}

func (g genThumbnail) getVideoDuration() (float64, error) {
	// get video duration
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-show-entries", "format=duration",
		"-of", "csv=p=0",
		g.file.FilePath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var duration float64
	fmt.Sscanf(string(output), "%f", &duration)

	return duration, nil
}

func (g genThumbnail) processVideo() error {
	dur, err := g.getVideoDuration()
	if err != nil {
		dur = 20
	}

	seekTime := max(5, dur*0.3)

	for _, ts := range g.ts {
		thumbnailDir := filepath.Join(filepath.Dir(g.file.FilePath), "thumbnails", ts.name)
		thumbnailPath := filepath.Join(thumbnailDir, g.file.StoredName+".jpg")

		if err := os.MkdirAll(thumbnailDir, 0755); err != nil {
			log.Printf("error creating %s thumbnail directory: %v\n", ts.name, err)
			continue
		}

		cmd := exec.Command("ffmpeg",
			"-i", g.file.FilePath,
			"-ss", fmt.Sprintf("%.1f", seekTime),
			"-vframes", "1",
			"-vf", fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2",
				ts.width, ts.height, ts.width, ts.height),
			"-y",
			thumbnailPath,
		)

		if err := cmd.Run(); err != nil {
			log.Printf("error creating %s thumbnail: %v\n", ts.name, err)
			continue
		}

		info, err := os.Stat(thumbnailPath)
		if err != nil {
			os.Remove(thumbnailPath)
			log.Printf("failed to stat thumbnail: %v\n", err)
			continue
		}

		relPath := "/static" + strings.TrimPrefix(thumbnailPath, g.conf.UploadPath())
		err = g.pg.InsertThumbnails(models.Thumbnail{
			FileID:       g.file.ID,
			SizeName:     ts.name,
			Width:        ts.width,
			Height:       ts.height,
			FilePath:     thumbnailPath,
			RelativePath: relPath,
			FileSize:     info.Size(),
		})

		if err != nil {
			os.Remove(thumbnailPath)
			log.Printf("failed to save thumbnail to db: %v\n", err)
			continue
		}
	}

	return nil
}

func (g genThumbnail) processImage() error {
	src, err := imaging.Open(g.file.FilePath)
	if err != nil {
		return fmt.Errorf("error opening: %v\n", err)
	}

	for _, ts := range g.ts {
		thumbnail := imaging.Fit(src, ts.width, ts.height, imaging.Lanczos)
		thumbnailDir := filepath.Join(filepath.Dir(g.file.FilePath), "thumbnails", ts.name)
		thumbnailPath := filepath.Join(thumbnailDir, g.file.StoredName)
		relPath := "/static" + strings.TrimPrefix(thumbnailPath, g.conf.UploadPath())

		if err := os.MkdirAll(thumbnailDir, 0755); err != nil {
			log.Printf("error creating %s thumbnail directory: %v\n", ts.name, err)
			continue
		}

		if err := imaging.Save(thumbnail, thumbnailPath); err != nil {
			log.Printf("error saving %s thumbnail: %v\n", ts.name, err)
			continue
		}

		info, err := os.Stat(thumbnailPath)
		if err != nil {
			log.Printf("failed to stat thumbnail: %v\n", err)
			os.Remove(thumbnailPath)
			continue
		}

		err = g.pg.InsertThumbnails(models.Thumbnail{
			FileID:       g.file.ID,
			SizeName:     ts.name,
			Width:        ts.width,
			Height:       ts.height,
			FilePath:     thumbnailPath,
			RelativePath: relPath,
			FileSize:     info.Size(),
		})
		if err != nil {
			log.Printf("failed to save thumbnail to db: %v\n", err)
			os.Remove(thumbnailPath)
			continue
		}
	}

	return nil
}

func (g genThumbnail) process() error {
	if strings.Contains(g.file.MimeType, "image") {
		return g.processImage()
	}

	if strings.Contains(g.file.MimeType, "video") {
		return g.processVideo()
	}

	return fmt.Errorf("gen thumbnail: unsupported type: %s", g.file.MimeType)
}
