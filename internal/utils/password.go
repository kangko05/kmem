package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func CheckPasswordHash(hash, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}
