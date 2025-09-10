package metadata

import (
	"context"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
)

var (
	testUid = "testUid"
	client  *firestore.Client
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	ctx := context.TODO()
	c, err := firestore.NewClient(ctx, os.Getenv("FIREBASE_PROJECT_ID"))
	if err != nil {
		log.Fatalf("failed to init client: %v", err)
	}
	client = c

	exitCode := m.Run()

	if err := c.Close(); err != nil {
		log.Printf("Error closing client: %v", err)
	}

	os.Exit(exitCode)
}

func TestMetadataAdd(t *testing.T) {
	cases := []struct {
		name          string
		setup         func(*firestore.Client)
		inputUid      string
		inputMetadata Metadata
		wantErr       bool
	}{
		{
			name:          "new metadata",
			setup:         func(c *firestore.Client) {},
			inputUid:      testUid,
			inputMetadata: FakeMetadata,
			wantErr:       false,
		},
		{
			name: "duplicate metadata",
			setup: func(c *firestore.Client) {
				client.Collection("projects").Doc(testUid).Set(
					context.TODO(),
					map[string]Metadata{FakeMetadata.Title: FakeMetadata},
				)
			},
			inputUid:      testUid,
			inputMetadata: FakeMetadata,
			wantErr:       true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() {
				_, _ = client.Collection("projects").Doc(testUid).Delete(context.TODO())
			})

			store := NewMetadataStore(client, os.Getenv("METADATA_COLLECTION"))
			err := store.MetadataAdd(tt.inputUid, tt.inputMetadata)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("%s: want error, got nil", tt.name)
				}
			} else {
				if err != nil {
					t.Fatalf(`%s: want nil, got %v`, tt.name, err)
				}
			}
		})
	}
}

func TestMetadataList(t *testing.T) {
	cases := []struct {
		name           string
		setup          func(client *firestore.Client)
		inputUid       string
		wantErr        bool
		wantListLength int
	}{
		{
			name: "List with one item",
			setup: func(client *firestore.Client) {
				_, _ = client.Collection("projects").Doc(testUid).Set(
					context.TODO(),
					map[string]Metadata{
						FakeMetadata.Title: FakeMetadata,
					},
				)
			},
			inputUid:       testUid,
			wantErr:        false,
			wantListLength: 1,
		},
		{
			name: "Empty list",
			setup: func(client *firestore.Client) {
				_, _ = client.Collection("projects").Doc(testUid).Set(
					context.TODO(),
					map[string]string{},
				)
			},
			inputUid:       testUid,
			wantErr:        false,
			wantListLength: 0,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() {
				_, _ = client.Collection("projects").Doc(testUid).Delete(context.TODO())
			})

			store := NewMetadataStore(client, os.Getenv("METADATA_COLLECTION"))
			got, err := store.MetadataList(tt.inputUid)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("%s: got nil, want error", tt.name)
				}
			} else {
				if err != nil {
					t.Fatalf(`%s: got %v, want nil`, tt.name, err)
				}
				if len(got) != tt.wantListLength {
					t.Fatalf(
						`%s: got list of length %d, want %d`,
						tt.name,
						len(got),
						tt.wantListLength,
					)
				}
			}
		})
	}
}

func TestMetadataUpdate(t *testing.T) {
	cases := []struct {
		name          string
		setup         func(*firestore.Client)
		inputUid      string
		inputMetadata Metadata
		wantErr       bool
	}{
		{
			name: "Update metadata title",
			setup: func(c *firestore.Client) {
				_, _ = client.Collection("projects").Doc(testUid).Set(
					context.TODO(),
					FakeMetadata,
				)
			},
			inputUid: testUid,
			inputMetadata: Metadata{
				Title:         "newTitle",
				Format:        FakeMetadata.Format,
				NotesMimeType: FakeMetadata.NotesMimeType,
				NotesSize:     FakeMetadata.NotesSize,
				Src:           FakeMetadata.Src,
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() {
				_, _ = client.Collection("projects").Doc(testUid).Delete(context.TODO())
			})

			store := NewMetadataStore(client, os.Getenv("METADATA_COLLECTION"))
			err := store.MetadataUpdate(tt.inputUid, tt.inputMetadata)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("%s: got nil, want error", tt.name)
				}
			} else {
				if err != nil {
					t.Fatalf(`%s: got %v, want nil`, tt.name, err)
				}
			}
		})
	}
}

func TestMetadataGet(t *testing.T) {
	cases := []struct {
		name         string
		setup        func(client *firestore.Client)
		inputUid     string
		inputId      string
		wantErr      bool
		wantMetadata Metadata
	}{
		{
			name: "Get test metadata",
			setup: func(client *firestore.Client) {
				_, _ = client.Collection("projects").Doc(testUid).Set(
					context.TODO(), map[string]Metadata{
						FakeMetadata.Title: FakeMetadata,
					},
				)
			},
			inputUid:     testUid,
			inputId:      FakeMetadata.Title,
			wantErr:      false,
			wantMetadata: FakeMetadata,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() {
				_, _ = client.Collection("projects").Doc(testUid).Delete(context.TODO())
			})

			store := NewMetadataStore(client, os.Getenv("METADATA_COLLECTION"))
			got, err := store.MetadataGet(tt.inputUid, tt.inputId)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("%s: got nil, want error", tt.name)
				}
			} else {
				if err != nil {
					t.Fatalf(`%s: got %v, want nil`, tt.name, err)
				}
				if got != tt.wantMetadata {
					t.Fatalf(`%s: got %q, want %q`, tt.name, got, tt.wantMetadata)
				}
			}
		})
	}
}

func TestMetadataRemove(t *testing.T) {
	cases := []struct {
		name     string
		setup    func(*firestore.Client)
		inputUid string
		inputId  string
		wantErr  bool
	}{
		{
			name: "Existent metadata",
			setup: func(c *firestore.Client) {
				_, _ = client.Collection("projects").Doc(testUid).Set(
					context.TODO(),
					FakeMetadata,
				)
			},
			inputUid: testUid,
			inputId:  testUid,
			wantErr:  false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() {
				_, _ = client.Collection("projects").Doc(testUid).Delete(context.TODO())
			})

			store := NewMetadataStore(client, os.Getenv("METADATA_COLLECTION"))
			err := store.MetadataRemove(testUid, FakeMetadata.Title)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("%s: got nil, want error", tt.name)
				}
			} else {
				if err != nil {
					t.Fatalf(`%s: got %v, want nil`, tt.name, err)
				}
			}
		})
	}
}

func TestMetadataRemoveAll(t *testing.T) {
	cases := []struct {
		name     string
		setup    func(*firestore.Client)
		inputUid string
		wantErr  bool
	}{
		{
			name: "Existent metadata",
			setup: func(c *firestore.Client) {
				_, _ = client.Collection("projects").Doc(testUid).Set(
					context.TODO(),
					FakeMetadata,
				)
			},
			inputUid: testUid,
			wantErr:  false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() {
				_, _ = client.Collection("projects").Doc(testUid).Delete(context.TODO())
			})

			store := NewMetadataStore(client, os.Getenv("METADATA_COLLECTION"))
			err := store.MetadataRemoveAll(testUid)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("%s: got nil, want error", tt.name)
				}
			} else {
				if err != nil {
					t.Fatalf(`%s: got %v, want nil`, tt.name, err)
				}
			}
		})
	}
}
