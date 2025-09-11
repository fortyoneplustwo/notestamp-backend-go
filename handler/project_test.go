package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"notestamp/media"
	"notestamp/metadata"
	"notestamp/notes"
	"notestamp/user"
	"testing"
)

func TestDeleteProject(t *testing.T) {
	cases := []struct {
		name           string
		uid            string
		id             string
		dr             metadata.MockMetadataStore
		mr             media.MockMediaStore
		nr             notes.MockNotesStore
		wantStatusCode int
	}{
		{
			name:           "successful delete",
			uid:            user.FakeUid,
			id:             metadata.FakeMetadata.Title,
			dr:             metadata.MockMetadataStore{Err: nil},
			mr:             media.MockMediaStore{Err: nil},
			nr:             notes.MockNotesStore{Err: nil},
			wantStatusCode: http.StatusNoContent,
		},
		{
			name:           "error deleting media",
			uid:            user.FakeUid,
			id:             metadata.FakeMetadata.Title,
			dr:             metadata.MockMetadataStore{Err: nil},
			mr:             media.MockMediaStore{Err: fakeError},
			nr:             notes.MockNotesStore{Err: nil},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(struct {
				Id string `json:"id"`
			}{
				Id: tt.id,
			})
			if err != nil {
				t.Fatalf("failed to contruct req body: %v", err)
			}

			r := httptest.NewRequest("DELETE", "/api/delete", bytes.NewBuffer(reqBody))
			ctx := user.NewContext(r.Context(), &user.AuthUser{Id: tt.uid})

			w := httptest.NewRecorder()

			handler := DeleteProject(tt.dr, tt.mr, tt.nr)
			handler.ServeHTTP(w, r.WithContext(ctx))

			if w.Result().StatusCode != tt.wantStatusCode {
				t.Fatalf("got %d, want %d", w.Result().StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestUpdateProject(t *testing.T) {
	cases := []struct {
		name           string
		uid            string
		id             string
		updatedNotes   string
		nu             notes.MockNotesStore
		wantStatusCode int
	}{
		{
			name:           "successful update",
			uid:            user.FakeUid,
			id:             user.FakeUid,
			updatedNotes:   notes.FakeNotes,
			nu:             notes.MockNotesStore{Err: nil},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "failed update",
			uid:            user.FakeUid,
			id:             user.FakeUid,
			updatedNotes:   notes.FakeNotes,
			nu:             notes.MockNotesStore{Err: fakeError},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(notes.Notes{
				Id:   tt.id,
				Data: []byte(tt.updatedNotes),
			})
			if err != nil {
				t.Fatalf("failed to contruct req body: %v", err)
			}

			r := httptest.NewRequest("PUT", "/api/update-notes", bytes.NewBuffer(reqBody))
			ctx := user.NewContext(r.Context(), &user.AuthUser{Id: tt.uid})

			w := httptest.NewRecorder()

			handler := UpdateProject(tt.nu)
			handler.ServeHTTP(w, r.WithContext(ctx))

			if w.Result().StatusCode != tt.wantStatusCode {
				t.Fatalf("got %d, want %d", w.Result().StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestStageWithoutMedia(t *testing.T) {
	cases := []struct {
		name           string
		uid            string
		metadata       metadata.Metadata
		staging        metadata.MockMetadataStore
		ns             notes.MockNotesStore
		wantStatusCode int
	}{
		{
			name:           "valid metadata",
			uid:            user.FakeUid,
			metadata:       metadata.FakeMetadata,
			staging:        metadata.MockMetadataStore{Err: nil},
			ns:             notes.MockNotesStore{Err: nil},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "invalid metadata (no src)",
			uid:  user.FakeUid,
			metadata: metadata.Metadata{
				Title:         metadata.FakeMetadata.Title,
				Format:        metadata.Youtube,
				Src:           "",
				NotesMimeType: metadata.FakeMetadata.NotesMimeType,
				NotesSize:     metadata.FakeMetadata.NotesSize,
			},
			staging:        metadata.MockMetadataStore{Err: nil},
			ns:             notes.MockNotesStore{Err: nil},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.metadata)
			if err != nil {
				t.Fatalf("could not construct req body: %v", err)
			}

			r := httptest.NewRequest("POST", "/api/save-without-media", bytes.NewBuffer(reqBody))
			ctx := user.NewContext(r.Context(), &user.AuthUser{Id: tt.uid})

			w := httptest.NewRecorder()

			handler := StageWithoutMedia(tt.staging, tt.ns)
			handler.ServeHTTP(w, r.WithContext(ctx))

			if w.Result().StatusCode != tt.wantStatusCode {
				t.Fatalf("got %d, want %d", w.Result().StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestStageWithMedia(t *testing.T) {
	cases := []struct {
		name           string
		uid            string
		metadata       metadata.Metadata
		staging        metadata.MockMetadataStore
		ms             media.MockMediaStore
		ns             notes.MockNotesStore
		wantStatusCode int
	}{
		{
			name: "valid metadata",
			uid:  user.FakeUid,
			metadata: metadata.Metadata{
				Title:         metadata.FakeMetadata.Title,
				Format:        metadata.Audio,
				Src:           "",
				NotesMimeType: metadata.FakeMetadata.NotesMimeType,
				NotesSize:     metadata.FakeMetadata.NotesSize,
				MediaMimeType: "audio/wav",
				MediaSize:     500,
			},
			staging:        metadata.MockMetadataStore{Err: nil},
			ms:             media.MockMediaStore{Err: nil},
			ns:             notes.MockNotesStore{Err: nil},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "missing media MIME type",
			uid:  user.FakeUid,
			metadata: metadata.Metadata{
				Title:         metadata.FakeMetadata.Title,
				Format:        metadata.Audio,
				NotesMimeType: metadata.FakeMetadata.NotesMimeType,
				NotesSize:     metadata.FakeMetadata.NotesSize,
				MediaSize:     500,
			},
			staging:        metadata.MockMetadataStore{Err: nil},
			ms:             media.MockMediaStore{Err: nil},
			ns:             notes.MockNotesStore{Err: nil},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "missing media size ",
			uid:  user.FakeUid,
			metadata: metadata.Metadata{
				Title:         metadata.FakeMetadata.Title,
				Format:        metadata.Audio,
				NotesMimeType: metadata.FakeMetadata.NotesMimeType,
				NotesSize:     metadata.FakeMetadata.NotesSize,
				MediaMimeType: "audio/wav",
			},
			staging:        metadata.MockMetadataStore{Err: nil},
			ms:             media.MockMediaStore{Err: nil},
			ns:             notes.MockNotesStore{Err: nil},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid media MIME type",
			uid:  user.FakeUid,
			metadata: metadata.Metadata{
				Title:         metadata.FakeMetadata.Title,
				Format:        metadata.Audio,
				NotesMimeType: metadata.FakeMetadata.NotesMimeType,
				NotesSize:     metadata.FakeMetadata.NotesSize,
				MediaSize:     500,
				MediaMimeType: "application/pdf",
			},
			staging:        metadata.MockMetadataStore{Err: nil},
			ms:             media.MockMediaStore{Err: nil},
			ns:             notes.MockNotesStore{Err: nil},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.metadata)
			if err != nil {
				t.Fatalf("could not construct req body: %v", err)
			}

			r := httptest.NewRequest("POST", "/api/save-without-media", bytes.NewBuffer(reqBody))
			ctx := user.NewContext(r.Context(), &user.AuthUser{Id: tt.uid})

			w := httptest.NewRecorder()

			handler := StageWithMedia(tt.staging, tt.ms, tt.ns)
			handler.ServeHTTP(w, r.WithContext(ctx))

			if w.Result().StatusCode != tt.wantStatusCode {
				t.Fatalf("got %d, want %d", w.Result().StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestPostSaveProject(t *testing.T) {
	cases := []struct {
		name           string
		uid            string
		id             string
		staging        metadata.MockMetadataStore
		ma             metadata.MockMetadataStore
		mc             media.MockMediaStore
		nc             notes.MockNotesStore
		wantStatusCode int
	}{
		{
			name: "project without media",
			uid: user.FakeUid,
			id: metadata.FakeMetadata.Title,
			staging: metadata.MockMetadataStore{Err: nil},
			ma: metadata.MockMetadataStore{Err: nil},
			mc: media.MockMediaStore{Err: nil},
			nc: notes.MockNotesStore{Err: nil},
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "missing notes",
			uid: user.FakeUid,
			id: metadata.FakeMetadata.Title,
			staging: metadata.MockMetadataStore{Err: nil},
			ma: metadata.MockMetadataStore{Err: nil},
			mc: media.MockMediaStore{Err: nil},
			nc: notes.MockNotesStore{Err: fakeError},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(struct {
				Id string `json:"id"`
			}{
				Id: tt.id,
			})
			if err != nil {
				t.Fatalf("failed to construct req body: %v", err)
			}

			r := httptest.NewRequest("POST", "/api/commit", bytes.NewReader(reqBody))
			ctx := user.NewContext(r.Context(), &user.AuthUser{Id: tt.uid})

			w := httptest.NewRecorder()

			handler := PostSaveProject(tt.staging, tt.ma, tt.mc, tt.nc)
			handler.ServeHTTP(w, r.WithContext(ctx))

			if w.Result().StatusCode != tt.wantStatusCode {
				var body bytes.Buffer
				io.Copy(&body, w.Result().Body)
				t.Fatalf("got %d: %s, want %d", w.Result().StatusCode, body.Bytes(), tt.wantStatusCode)
			}
		})
	}
}
