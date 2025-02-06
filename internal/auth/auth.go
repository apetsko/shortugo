package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/apetsko/shortugo/internal/utils"
	"github.com/gorilla/securecookie"
)

func securedCookie(secret string) *securecookie.SecureCookie {
	secretLen := 32
	id := utils.GenerateID(secret, secretLen)
	sharedSecret := []byte(id)
	return securecookie.New(sharedSecret, sharedSecret)
}

func SetUserIDCookie(w http.ResponseWriter, secret string) (userID string, err error) {
	sc := securedCookie(secret)

	userIDLen := 8
	userID, err = utils.GenerateUserID(userIDLen)
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

func UserIDFromCookie(r *http.Request, secret string) (string, error) {
	cookie, err := r.Cookie("shortugo")
	if err != nil {
		err = fmt.Errorf("error getting userid cookie: %w", err)
		return "", err
	}

	userID, err := GetUserIDFromCookie(cookie, secret)
	if err != nil {
		err = fmt.Errorf("error getting userid from cookie: %w", err)
		return "", err
	}
	return userID, nil
}
