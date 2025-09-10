package middleware

import (
	"net/http"
	"net/http/httptest"
	"notestamp/auth"
	jwt "notestamp/auth/jwt"
	"notestamp/user"
	"testing"
	"time"
)

type mockHandler func()

func (h mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func TestAuthMiddleware(t *testing.T) {
	cases := []struct {
		name           string
		uid            string
		revokedStore   auth.MockRevokedStore
		includeCookies bool
		accExpTime     time.Time
		refExpTime     time.Time
		wantStatus     int
	}{
		{
			name:           "authorized",
			uid:            user.FakeUid,
			revokedStore:   auth.MockRevokedStore{IsRevoked: false},
			includeCookies: true,
			accExpTime:     time.Now().Add(time.Hour),
			refExpTime:     time.Now().Add(time.Hour),
			wantStatus:     http.StatusOK,
		},
		{
			name:           "authorized (expired access token)",
			uid:            user.FakeUid,
			revokedStore:   auth.MockRevokedStore{IsRevoked: false},
			includeCookies: true,
			accExpTime:     time.Now().Add(-time.Hour),
			refExpTime:     time.Now().Add(time.Hour),
			wantStatus:     http.StatusOK,
		},
		{
			name:           "unauthorized (missing tokens)",
			uid:            user.FakeUid,
			revokedStore:   auth.MockRevokedStore{IsRevoked: false},
			includeCookies: false,
			accExpTime:     time.Now().Add(time.Hour),
			refExpTime:     time.Now().Add(time.Hour),
			wantStatus:     http.StatusUnauthorized,
		},
		{
			name:           "unauthorized (expired refresh token)",
			uid:            user.FakeUid,
			revokedStore:   auth.MockRevokedStore{IsRevoked: false},
			includeCookies: true,
			accExpTime:     time.Now().Add(-time.Hour),
			refExpTime:     time.Now().Add(-time.Hour),
			wantStatus:     http.StatusUnauthorized,
		},
		{
			name:           "unauthorized (refresh token revoked)",
			uid:            user.FakeUid,
			revokedStore:   auth.MockRevokedStore{IsRevoked: true},
			includeCookies: true,
			accExpTime:     time.Now().Add(-time.Hour),
			refExpTime:     time.Now().Add(time.Hour),
			wantStatus:     http.StatusUnauthorized,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var h mockHandler = func() {}
			mw := Authenticate(h, tt.revokedStore)

			r := httptest.NewRequest("GET", "/api/", nil)
			w := httptest.NewRecorder()

			if tt.includeCookies {
				acc, ref, err := jwt.GenerateTokens(user.FakeUid, tt.accExpTime, tt.refExpTime)
				if err != nil {
					t.Fatalf("could not generate tokens: %v", err)
				}
				r.AddCookie(&http.Cookie{Name: "access-token", Value: acc})
				r.AddCookie(&http.Cookie{Name: "refresh-token", Value: ref})
			}

			mw.ServeHTTP(w, r)

			if w.Result().StatusCode != tt.wantStatus {
				t.Fatalf("got http status: %d, want %d", w.Result().StatusCode, tt.wantStatus)
			}
		})
	}
}
