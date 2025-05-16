package models

import "time"

type FileMetadata struct {
	Filename    string    `json:"filename"`
	Size        int64     `json:"size"`
	ContentType string    `json:"contentType"`
	UploadedBy  string    `json:"uploadedBy"`
	UploadedAt  time.Time `json:"uploadedAt"`
}
