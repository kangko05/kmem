package models

import (
	"fmt"
	"time"
)

type File struct {
	ID           int       `json:"id" db:"id"`
	Hash         string    `json:"hash" db:"hash"` // md5
	Username     string    `json:"username" db:"username"`
	OriginalName string    `json:"originalName" db:"original_name"`
	StoredName   string    `json:"storedName" db:"stored_name"`
	FilePath     string    `json:"filePath" db:"file_path"`
	FileSize     int64     `json:"fileSize" db:"file_size"`
	MimeType     string    `json:"mimeType" db:"mime_type"`
	UploadedAt   time.Time `json:"uploadedAt" db:"uploaded_at"`
}

func (f *File) IsImage() bool {
	switch f.MimeType {
	case "image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}

func (f *File) IsVideo() bool {
	switch f.MimeType {
	case "video/mp4", "video/avi", "video/mov", "video/mkv", "video/webm":
		return true
	default:
		return false
	}
}

func (f *File) GetFileType() string {
	if f.IsImage() {
		return "image"
	}
	if f.IsVideo() {
		return "video"
	}
	return "other"
}

func (f *File) GetReadableSize() string {
	const unit = 1024
	if f.FileSize < unit {
		return fmt.Sprintf("%d B", f.FileSize)
	}

	div, exp := int64(unit), 0
	for n := f.FileSize / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(f.FileSize)/float64(div), "KMGTPE"[exp])
}

// DTO ========================================================================

type FileUploadRequest struct {
	Filename string `json:"filename" form:"filename"`
}

type FileUploadResponse struct {
	ID           int    `json:"id"`
	OriginalName string `json:"originalName"`
	StoredName   string `json:"storedName"`
	FileSize     int64  `json:"fileSize"`
	MimeType     string `json:"mimeType"`
	FileType     string `json:"fileType"` // "image", "video", "other"
}

// type FileListResponse struct {
// 	Files      []File `json:"files"`
// 	TotalCount int    `json:"totalCount"`
// 	Page       int    `json:"page"`
// 	Limit      int    `json:"limit"`
// }
