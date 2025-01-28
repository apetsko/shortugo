package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
)

func GenerateID(s string, length int) (id string) {
	hash := sha256.Sum256([]byte(s))
	id = base64.RawURLEncoding.EncodeToString(hash[:length])[:length]
	return
}

func GenerateUserID(length int) (id string, err error) {
	r := make([]byte, length)

	_, err = rand.Read(r)
	if err != nil {
		err = fmt.Errorf("failed to generate random User ID: %w", err)
		return "", err
	}
	id = hex.EncodeToString(r)

	return id, nil
}

func securedCookie(secret string) *securecookie.SecureCookie {
	secretLen := 32
	id := GenerateID(secret, secretLen)
	sharedSecred := []byte(id)
	return securecookie.New(sharedSecred, sharedSecred)
}

func SetUserCookie(w http.ResponseWriter, secret string) (userID string, err error) {
	sc := securedCookie(secret)

	userIDLen := 8
	userID, err = GenerateUserID(userIDLen)
	if err != nil {
		return "", err
	}

	encoded, err := sc.Encode("shortugo", userID)
	if err != nil {
		err = fmt.Errorf("error encoding userid cookie: %v", err)
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "shortugo",
		Value:    encoded,
		HttpOnly: true,
		Path:     "/",
	})
	return userID, nil
}

func GetUserIDFromCookie(cookie *http.Cookie, secret string) (string, error) {
	var userID string
	sc := securedCookie(secret)
	if err := sc.Decode("shortugo", cookie.Value, &userID); err != nil {
		err = fmt.Errorf("error decoding user cookie: %w", err)
		return "", err
	}

	if userID == "" {
		return "", errors.New("userid not found in cookie")
	}
	return userID, nil
}

func CookieUserID(r *http.Request, secret string) (string, error) {
	cookie, err := r.Cookie("shortugo")
	if err != nil {
		return "", err
	}

	userID, err := GetUserIDFromCookie(cookie, secret)
	if err != nil {
		return "", err
	}
	return userID, nil
}
