package models

import "time"

type Thumbnail struct {
	ID           int       `json:"id" db:"id"`
	FileID       int       `json:"fileId" db:"file_id"`
	SizeName     string    `json:"sizeName" db:"size_name"`
	Width        int       `json:"width" db:"width"`
	Height       int       `json:"height" db:"height"`
	FilePath     string    `json:"filePath" db:"file_path"`
	FileSize     int64     `json:"fileSize" db:"file_size"`
	RelativePath string    `json:"relativePath" db:"relative_path"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

type ThumbnailResponse struct {
	SizeName string `json:"sizeName,omitempty"`
	FilePath string `json:"filePath,omitempty"`
}
