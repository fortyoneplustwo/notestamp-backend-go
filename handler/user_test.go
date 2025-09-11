package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"notestamp/auth"
	jwt "notestamp/auth/jwt"
	store "notestamp/handler/interfaces"
	"notestamp/media"
	"notestamp/metadata"
	"notestamp/notes"
	"notestamp/user"
	"testing"
	"time"
)

var (
	fakeError = errors.New("fakeError")
)

func TestRegisterUser(t *testing.T) {
	cases := []struct {
		name       string
		userStore  store.UserAdder
		payload    user.Credentials
		wantStatus int
	}{
		{
			name:       "new user",
			userStore:  user.MockUserStore{Err: nil},
			payload:    user.FakeCredentials,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "registered user",
			userStore:  user.MockUserStore{Err: fakeError},
			payload:    user.FakeCredentials,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			handler := RegisterUser(tt.userStore)

			body, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatalf("could not encode payload to json")
			}

			r := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			ctx := context.WithValue(r.Context(), "uid", user.FakeUid)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r.WithContext(ctx))

			if w.Code != tt.wantStatus {
				t.Fatalf("got http status: %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestDeregisterUser(t *testing.T) {
	cases := []struct {
		name              string
		userRemover       user.MockUserStore
		metadataRemover   metadata.MockMetadataStore
		revokedTokenAdder auth.MockRevokedStore
		notesRemover      notes.MockNotesStore
		mediaRemover      media.MockMediaStore
		isAuthorized      bool
		wantStatus        int
	}{
		{
			name:              "registered user",
			userRemover:       user.MockUserStore{Err: nil},
			metadataRemover:   metadata.MockMetadataStore{Err: nil},
			revokedTokenAdder: auth.MockRevokedStore{Err: nil},
			notesRemover:      notes.MockNotesStore{Err: nil},
			mediaRemover:      media.MockMediaStore{Err: nil},
			isAuthorized:      true,
			wantStatus:        http.StatusNoContent,
		},
		{
			name:              "non-existent user",
			userRemover:       user.MockUserStore{Err: fakeError},
			metadataRemover:   metadata.MockMetadataStore{Err: fakeError},
			revokedTokenAdder: auth.MockRevokedStore{Err: fakeError},
			notesRemover:      notes.MockNotesStore{Err: nil},
			mediaRemover:      media.MockMediaStore{Err: nil},
			isAuthorized:      true,
			wantStatus:        http.StatusInternalServerError,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			handler := DeregisterUser(
				tt.userRemover,
				tt.metadataRemover,
				tt.mediaRemover,
				tt.notesRemover,
				tt.revokedTokenAdder,
			)

			r := httptest.NewRequest("POST", "/register", nil)
			if tt.isAuthorized {
				accExp := time.Now().Add(time.Hour)
				refExp := time.Now().Add(time.Hour)
				acc, ref, err := jwt.GenerateTokens(user.FakeUid, accExp, refExp)
				if err != nil {
					t.Fatalf("could not generate tokens")
				}
				r.AddCookie(&http.Cookie{Name: "access-token", Value: acc})
				r.AddCookie(&http.Cookie{Name: "refresh-token", Value: ref})
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			if w.Code != tt.wantStatus {
				t.Fatalf("got http status: %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	cases := []struct {
		name              string
		userGetter        user.MockUserStore
		revokedTokenAdder auth.MockRevokedStore
		payload           user.Credentials
		isAuthorized      bool
		wantStatus        int
	}{
		{
			name:              "registered authorized user",
			userGetter:        user.MockUserStore{Err: nil, User: user.FakeUser},
			revokedTokenAdder: auth.MockRevokedStore{Err: nil},
			payload:           user.FakeCredentials,
			isAuthorized:      true,
			wantStatus:        http.StatusOK,
		},
		{
			name:              "unregistered authorized user",
			userGetter:        user.MockUserStore{Err: fakeError},
			revokedTokenAdder: auth.MockRevokedStore{Err: nil},
			payload:           user.FakeCredentials,
			isAuthorized:      true,
			wantStatus:        http.StatusUnauthorized,
		},
		{
			name: "wrong password",
			userGetter: user.MockUserStore{
				Err:  fakeError,
				User: user.FakeUser,
			},
			revokedTokenAdder: auth.MockRevokedStore{Err: nil},
			payload: user.Credentials{
				Email:    user.FakeCredentials.Email,
				Password: "fakeWrongPassword",
			},
			isAuthorized: true,
			wantStatus:   http.StatusUnauthorized,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := auth.CreateHash(tt.userGetter.User.Password)
			if err != nil {
				t.Fatalf("could not hash password: %v", err)
			}
			tt.userGetter.User.Password = hash

			handler := LoginUser(tt.userGetter, tt.revokedTokenAdder)

			payload, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatalf("could not contruct request body: %v", err)
			}

			r := httptest.NewRequest("PUT", "/login", bytes.NewBuffer(payload))
			if tt.isAuthorized {
				accExp := time.Now().Add(time.Hour)
				refExp := time.Now().Add(time.Hour)
				acc, ref, err := jwt.GenerateTokens(user.FakeUid, accExp, refExp)
				if err != nil {
					t.Fatalf("could not generate tokens")
				}
				r.AddCookie(&http.Cookie{Name: "access-token", Value: acc})
				r.AddCookie(&http.Cookie{Name: "refresh-token", Value: ref})
			}

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			if w.Code != tt.wantStatus {
				msg, _ := io.ReadAll(w.Body)
				t.Fatalf("got http status: %d %s, want %d", w.Code, msg, tt.wantStatus)
			}
		})
	}
}

func TestLogoutUser(t *testing.T) {
	cases := []struct {
		name              string
		revokedTokenAdder auth.MockRevokedStore
		isAuthorized      bool
		wantStatus        int
	}{
		{
			name:              "authorized user",
			revokedTokenAdder: auth.MockRevokedStore{Err: nil},
			isAuthorized:      true,
			wantStatus:        http.StatusOK,
		},
		{
			name:              "unauthorized user",
			revokedTokenAdder: auth.MockRevokedStore{Err: nil},
			isAuthorized:      false,
			wantStatus:        http.StatusUnauthorized,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			handler := LogoutUser(tt.revokedTokenAdder)

			r := httptest.NewRequest("PUT", "/logout", nil)

			if tt.isAuthorized {
				accExp := time.Now().Add(time.Hour)
				refExp := time.Now().Add(time.Hour)
				acc, ref, err := jwt.GenerateTokens(user.FakeUid, accExp, refExp)
				if err != nil {
					t.Fatalf("could not generate tokens")
				}
				r.AddCookie(&http.Cookie{Name: "access-token", Value: acc})
				r.AddCookie(&http.Cookie{Name: "refresh-token", Value: ref})
			}

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			if w.Result().StatusCode != tt.wantStatus {
				t.Fatalf(
					"got http status %d, want %d",
					w.Result().StatusCode,
					tt.wantStatus,
				)
			}
		})
	}
}
