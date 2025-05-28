package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// md5
func HashFile(file []byte) (string, error) {
	hash := md5.New()
	_, err := hash.Write(file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
