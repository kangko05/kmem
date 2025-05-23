package models

import (
	"time"
)

type FileMetadata struct {
	Hash        string    `json:"hash"` // for unique id, md5
	Filename    string    `json:"filename"`
	ContentType string    `json:"contentType"`
	StoredPath  string    `json:"storedPath"`
	ArchivePath string    `json:"archivePath"`
	UploadedBy  string    `json:"uploadedBy"`
	UploadedAt  time.Time `json:"uploadedAt"`
	Size        int64     `json:"size"`
}

type MetadataPart struct {
	Filename    string    `json:"filename"`
	Storedpath  string    `json:"storedpath"`
	ContentType string    `json:"contenttype"`
	UploadedAt  time.Time `json:"uploadedat"`
	Size        int64     `json:"size"`
}
