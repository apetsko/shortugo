package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/utils"
)

func WithUserID(ctx context.Context, userid string) context.Context {
	userIDkey := models.UserID("userid")
	return context.WithValue(ctx, userIDkey, userid)
}

func AuthMiddleware(secret string, logger *logging.ZapLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				if r.URL.Path == "/api/user/urls" {
					cookie, err := r.Cookie("shortugo")
					if err != nil {
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
					userID, err := utils.GetUserIDFromCookie(cookie, secret)
					if err != nil {
						err = fmt.Errorf("error getting userid from cookie: %w", err)
						logger.Error(err.Error())
						w.WriteHeader(http.StatusUnauthorized)
						return
					}

					r = r.WithContext(WithUserID(r.Context(), userID))
				}
			case http.MethodPost:
				userID, err := utils.CookieUserID(r, secret)
				if err != nil {
					userID, err = utils.SetUserCookie(w, secret)
					if err != nil {
						logger.Error(err.Error())
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
				r = r.WithContext(WithUserID(r.Context(), userID))
			}
			next.ServeHTTP(w, r)
		})
	}
}
