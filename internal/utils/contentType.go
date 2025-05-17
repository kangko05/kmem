package utils

import (
	"path/filepath"
	"strings"
)

func GetContentType(filename string) ContentType {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return IMAGE

	case ".mp4", ".avi", ".mov", ".wmv", ".mkv":
		return VIDEO

	case ".mp3", ".wav", ".ogg", ".flac", ".aac", ".m4a":
		return AUDIO

	case ".zip", ".tar", ".gz", ".rar", ".7z":
		return ZIP

	case ".txt", ".pdf", ".doc", ".docx":
		return TEXT

	default:
		return OCTET_STREAM
	}
}
