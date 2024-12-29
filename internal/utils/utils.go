package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)


func Generate(URL string) (ID string) {
	hash := sha256.Sum256([]byte(URL))
	ID = base64.RawURLEncoding.EncodeToString(hash[:6])
	return
}

func FullURL(baseURL string, ID string) string {
	return fmt.Sprintf("%s/%s", baseURL, ID)
}
