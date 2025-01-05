package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

func Generate(url string) (id string) {
	hash := sha256.Sum256([]byte(url))
	id = base64.RawURLEncoding.EncodeToString(hash[:6])
	return
}
