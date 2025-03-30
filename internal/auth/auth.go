package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/apetsko/shortugo/internal/utils"
	"github.com/gorilla/securecookie"
)

type Authenticator interface {
	CookieGetUserID(r *http.Request, secret string) (string, error)
	CookieSetUserID(w http.ResponseWriter, secret string) (userID string, err error)
}

type Auth struct{}

var ErrNoUserIDFound = errors.New("no user ID found in cookie")

func (a *Auth) CookieGetUserID(r *http.Request, secret string) (string, error) {
	cookie, err := r.Cookie("shortugo")
	if err != nil {
		return "", http.ErrNoCookie
	}

	sc := securedCookie(secret)

	var userID string

	err = sc.Decode("shortugo", cookie.Value, &userID)
	if err != nil || userID == "" {
		err = fmt.Errorf("%w: %w", ErrNoUserIDFound, err)
		return "", err
	}

	return userID, nil
}

func securedCookie(secret string) *securecookie.SecureCookie {
	secretLen := 32
	id := utils.GenerateID(secret, secretLen)
	sharedSecret := []byte(id)
	return securecookie.New(sharedSecret, sharedSecret)
}

func (a *Auth) CookieSetUserID(w http.ResponseWriter, secret string) (userID string, err error) {
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
