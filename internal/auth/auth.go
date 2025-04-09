package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/apetsko/shortugo/internal/utils"
	"github.com/gorilla/securecookie"
)

// Authenticator defines methods for handling authentication via cookies.
type Authenticator interface {
	// CookieGetUserID extracts the user ID from the request's cookie.
	// Returns the userID or an error if the cookie is invalid or missing.
	CookieGetUserID(r *http.Request, secret string) (string, error)

	// CookieSetUserID sets the user ID in the response's cookie.
	// Returns the set userID or an error if the setting fails.
	CookieSetUserID(w http.ResponseWriter, secret string) (userID string, err error)
}

// Auth is a struct that handles user authentication-related operations.
// It includes methods for setting and getting user IDs via cookies.
type Auth struct{}

var errNoUserIDFound = errors.New("no user ID found in cookie")

// CookieGetUserID retrieves the user ID from the cookie in the request.
// It decodes the cookie value and returns the user ID if found, or an error if not.
func (a *Auth) CookieGetUserID(r *http.Request, secret string) (string, error) {
	// Get the "shortugo" cookie from the request
	cookie, err := r.Cookie("shortugo")
	if err != nil {
		// Return an error if the cookie is not found
		return "", http.ErrNoCookie
	}

	// Create a secured cookie with the given secret
	sc := securedCookie(secret)

	var userID string

	// Decode the user ID from the cookie
	err = sc.Decode("shortugo", cookie.Value, &userID)
	if err != nil || userID == "" {
		// Return an error if decoding fails or user ID is empty
		err = fmt.Errorf("%w: %w", errNoUserIDFound, err)
		return "", err
	}

	// Return the decoded user ID
	return userID, nil
}

func securedCookie(secret string) *securecookie.SecureCookie {
	secretLen := 32
	id := utils.GenerateID(secret, secretLen)
	sharedSecret := []byte(id)
	return securecookie.New(sharedSecret, sharedSecret)
}

// CookieSetUserID sets the user ID in a secured cookie in the response.
// It generates a new user ID, encodes it, and sets it as a cookie in the response.
func (a *Auth) CookieSetUserID(w http.ResponseWriter, secret string) (userID string, err error) {
	// Create a secured cookie with the given secret
	sc := securedCookie(secret)

	// Generate a new user ID with a specified length
	userIDLen := 8
	userID, err = utils.GenerateUserID(userIDLen)
	if err != nil {
		return "", err
	}

	// Encode the user ID into the cookie
	encoded, err := sc.Encode("shortugo", userID)
	if err != nil {
		// Return an error if encoding fails
		err = fmt.Errorf("error encoding userid cookie: %v", err)
		return "", err
	}

	// Set the cookie in the response with the encoded user ID
	http.SetCookie(w, &http.Cookie{
		Name:     "shortugo",
		Value:    encoded,
		HttpOnly: true,
		Path:     "/",
	})
	return userID, nil
}
