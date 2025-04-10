package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apetsko/shortugo/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth_CookieSetUserID(t *testing.T) {
	auth := &Auth{}
	secret := "test_secret"

	testCases := []struct {
		wantErr error
		name    string
		secret  string
	}{
		{
			name:    "Valid cookie set",
			secret:  secret,
			wantErr: nil,
		},
		{
			name:    "Empty secret",
			secret:  "",
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			userID, err := auth.CookieSetUserID(w, tc.secret)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.Len(t, userID, 16)

				resp := w.Result()
				defer func() {
					require.NoError(t, resp.Body.Close())
				}()
				cookie := resp.Cookies()
				require.Len(t, cookie, 1)
				assert.Equal(t, "shortugo", cookie[0].Name)
			}
		})
	}
}

func TestAuth_CookieGetUserID(t *testing.T) {
	auth := &Auth{}
	secret := "test_secret"

	testCases := []struct {
		wantErr     error
		name        string
		secret      string
		setupCookie bool
	}{
		{
			name:        "Valid user ID retrieved",
			setupCookie: true,
			secret:      secret,
			wantErr:     nil,
		},
		{
			name:        "No cookie present",
			setupCookie: false,
			secret:      secret,
			wantErr:     http.ErrNoCookie,
		},
		{
			name:        "Invalid secret",
			setupCookie: true,
			secret:      "wrong_secret",
			wantErr:     errNoUserIDFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			if tc.setupCookie {
				userID, err := auth.CookieSetUserID(w, secret)
				require.NoError(t, err)

				resp := w.Result()
				defer func() {
					require.NoError(t, resp.Body.Close())
				}()
				cookies := resp.Cookies()
				require.Len(t, cookies, 1)

				r.AddCookie(cookies[0])

				gotUserID, err := auth.CookieGetUserID(r, tc.secret)

				if tc.wantErr != nil {
					assert.ErrorIs(t, err, tc.wantErr)
				} else {
					require.NoError(t, err)
					assert.Equal(t, userID, gotUserID)
				}
			} else {
				_, err := auth.CookieGetUserID(r, tc.secret)
				assert.ErrorIs(t, err, tc.wantErr)
			}
		})
	}
}

func TestSecuredCookie(t *testing.T) {
	testCases := []struct {
		name   string
		secret string
	}{
		{"Valid secret", "test_secret"},
		{"Empty secret", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sc := securedCookie(tc.secret)
			assert.NotNil(t, sc)
		})
	}
}

func TestGenerateSecureCookieSecret(t *testing.T) {
	secret := "test_secret"

	id1 := utils.GenerateID(secret, 32)
	id2 := utils.GenerateID(secret, 32)

	assert.Equal(t, id1, id2)
	assert.Len(t, id1, 32)
}
