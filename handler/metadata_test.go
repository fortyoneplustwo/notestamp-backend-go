package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"notestamp/metadata"
	"notestamp/user"
	"testing"
)

func TestListMetadata(t *testing.T) {
	cases := []struct {
		name       string
		uid        string
		wantErr    bool
		wantStatus int
		wantList   []metadata.Metadata
	}{
		{
			name:       "List of length 1",
			uid:        user.FakeUid,
			wantErr:    false,
			wantStatus: http.StatusOK,
			wantList:   []metadata.Metadata{metadata.FakeMetadata},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			handler := ListMetadata(metadata.MockMetadataStore{})

			req, err := http.NewRequest("GET", "/api/list", nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := user.NewContext(req.Context(), &user.AuthUser{Id: tt.uid})

			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req.WithContext(ctx))

			if resp.Code != tt.wantStatus {
				body := resp.Body
				t.Fatalf("got http status: %d: %s, want %d", resp.Code, body, tt.wantStatus)
			}

			if !tt.wantErr {
				payload := struct {
					Projects []metadata.Metadata `json:"projects"`
				}{}
				json.Unmarshal(resp.Body.Bytes(), &payload)
				list := payload.Projects
				if len(list) != 1 {
					t.Fatalf("got list of length: %d, want %d", len(list), len(tt.wantList))
				}
				if list[0] != tt.wantList[0] {
					t.Fatalf("got list[0] == %q, want %q", list[0], tt.wantList[0])
				}
			}
		})
	}
}
