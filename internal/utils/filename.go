package utils

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

func DecodeFilename(encodedFilename string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedFilename)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	decodedFilename, err := url.QueryUnescape(string(decodedBytes))
	if err != nil {
		return "", fmt.Errorf("failed to decode URL: %v", err)
	}

	return decodedFilename, nil
}

func EncodeFilename(filename string) string {
	encoded := url.QueryEscape(filename)

	return base64.StdEncoding.EncodeToString([]byte(encoded))
}

func ValidateFilename(filename string) error {
	if len(filename) == 0 {
		return fmt.Errorf("filename is empty")
	}

	if len(filename) > 255 {
		return fmt.Errorf("filename too long (max 255 characters)")
	}

	dangerousChars := []string{"..", "/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range dangerousChars {
		if strings.Contains(filename, char) {
			return fmt.Errorf("filename contains invalid character: %s", char)
		}
	}

	if strings.HasPrefix(filename, ".") {
		return fmt.Errorf("hidden files not allowed")
	}

	return nil
}

func GetAllowedMimeType(filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	allowedTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",

		".mp4":  "video/mp4",
		".avi":  "video/avi",
		".mov":  "video/quicktime",
		".mkv":  "video/x-matroska",
		".webm": "video/webm",
	}

	mimeType, allowed := allowedTypes[ext]
	if !allowed {
		return "", fmt.Errorf("file type not allowed: %s", ext)
	}

	return mimeType, nil
}

func ProcessFilename(encodedFilename string) (originalName, safeName, mimeType string, err error) {
	originalName, err = DecodeFilename(encodedFilename)
	if err != nil {
		return "", "", "", fmt.Errorf("decode error: %v", err)
	}

	if err = ValidateFilename(originalName); err != nil {
		return "", "", "", fmt.Errorf("validation error: %v", err)
	}

	mimeType, err = GetAllowedMimeType(originalName)
	if err != nil {
		return "", "", "", fmt.Errorf("mime type error: %v", err)
	}

	ext := filepath.Ext(originalName)
	safeName = GenerateUniqueFilename() + ext

	return originalName, safeName, mimeType, nil
}

func GenerateUniqueFilename() string {
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), rand.Intn(10000))
}
