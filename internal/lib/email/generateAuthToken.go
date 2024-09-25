package email

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomAuthString() (string, error) {
	b := make([]byte, 30)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
