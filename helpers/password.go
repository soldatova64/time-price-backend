package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

func HashPassword(password string) string {
	return "hashed_" + password
}

func GenerateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
