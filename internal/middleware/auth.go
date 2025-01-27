package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/utils"
	"github.com/gorilla/securecookie"
)

type userCookie struct {
	UserID string
}

func setUserCookie(w http.ResponseWriter, sc *securecookie.SecureCookie, userID string) error {
	encoded, err := sc.Encode("shortugo", userCookie{UserID: userID})
	if err != nil {
		err = fmt.Errorf("error encoding user cookie: %v", err)
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "shortugo",
		Value:    encoded,
		HttpOnly: true,
		Path:     "/",
	})
	return nil
}

func getUserIDFromCookie(cookie *http.Cookie, sc *securecookie.SecureCookie) (string, error) {
	uc := new(userCookie)

	if err := sc.Decode("shortugo", cookie.Value, &uc); err != nil {
		err = fmt.Errorf("error decoding user cookie: %w", err)
		return "", err
	}

	userID := uc.UserID
	if userID == "" {
		return "", errors.New("user_id not found in cookie")
	}
	return userID, nil
}

func WithUserID(ctx context.Context, uid string) context.Context {
	userid := models.UserID("UserID")
	return context.WithValue(ctx, userid, uid)
}

func AuthMiddleware(secret string, logger *logging.ZapLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secretLen := 32
			id := utils.GenerateID(secret, secretLen)
			sharedSecred := []byte(id)
			sc := securecookie.New(sharedSecred, sharedSecred)

			cookie, err := r.Cookie("shortugo")
			if err != nil {
				userIDLen := 10
				newUserID, err := utils.GenerateUserID(userIDLen)
				if err != nil {
					logger.Error(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if err = setUserCookie(w, sc, newUserID); err != nil {
					logger.Error(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				r = r.WithContext(WithUserID(r.Context(), newUserID))
				next.ServeHTTP(w, r)
				return
			}

			userID, err := getUserIDFromCookie(cookie, sc)
			if err != nil || userID == "" {
				err = fmt.Errorf("error getting userid from cookie: %w", err)
				logger.Error(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			r = r.WithContext(WithUserID(r.Context(), userID))
			next.ServeHTTP(w, r)
		})
	}
}
