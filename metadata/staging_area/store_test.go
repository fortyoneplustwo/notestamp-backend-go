package metadata

// NOTE: Must start redis server before running this test: $ redis-server

import (
	"encoding/json"
	"fmt"
	"log"
	"notestamp/metadata"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

var (
	client *redis.Client
)

func cleanupStagingArea(c *redis.Client) {
	err := c.FlushAll().Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to flush all databases: %v", err)
	}
}

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	exitCode := m.Run()

	err := client.Close()
	if err != nil {
		log.Printf("failed to close client: %v", err)
	}

	os.Exit(exitCode)
}

func TestMedataAdd(t *testing.T) {
	staging := NewStagingArea(client, time.Hour)

	cases := []struct {
		name     string
		setup    func()
		uid      string
		metadata metadata.Metadata
		wantErr  bool
	}{
		{
			name:     "new metadata",
			setup:    func() {},
			uid:      metadata.FakeUid,
			metadata: metadata.FakeMetadata,
			wantErr:  false,
		},
		{
			name: "overwrite metadata",
			setup: func() {
				key := fmt.Sprintf("%s/%s", metadata.FakeUid, metadata.FakeMetadata.Title)
				val, err := json.Marshal(metadata.FakeMetadata)
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				err = staging.client.Set(key, val, staging.expiresIn).Err()
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			uid:      metadata.FakeUid,
			metadata: metadata.FakeMetadata,
			wantErr:  false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			t.Cleanup(func() { cleanupStagingArea(client) })

			err := staging.MetadataAdd(tt.uid, tt.metadata)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("got nil, want err")
				}
			} else {
				if err != nil {
					t.Fatalf("got error: %v, want nil", err)
				}
			}
		})
	}
}

func TestMetadataGet(t *testing.T) {
	staging := NewStagingArea(client, time.Hour)

	cases := []struct {
		name         string
		setup        func()
		uid          string
		id           string
		wantErr      bool
		wantMetadata metadata.Metadata
	}{
		{
			name: "existent metadata",
			setup: func() {
				key := fmt.Sprintf("%s/%s", metadata.FakeUid, metadata.FakeMetadata.Title)
				val, err := json.Marshal(metadata.FakeMetadata)
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				err = client.Set(key, val, staging.expiresIn).Err()
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			uid:          metadata.FakeUid,
			id:           metadata.FakeMetadata.Title,
			wantErr:      false,
			wantMetadata: metadata.FakeMetadata,
		},
		{
			name:         "non-existent metadata",
			setup:        func() {},
			uid:          metadata.FakeUid,
			id:           metadata.FakeMetadata.Title,
			wantErr:      true,
			wantMetadata: metadata.Metadata{},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			t.Cleanup(func() { cleanupStagingArea(client) })

			m, err := staging.MetadataGet(tt.uid, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("got nil, want err")
				}
			} else {
				if err != nil {
					t.Fatalf("got err: %v, want nil", err)
				}
				if m != tt.wantMetadata {
					t.Fatalf("got %q want %q", m, tt.wantMetadata)
				}
			}
		})
	}
}

func TestMetadataRemove(t *testing.T) {
	staging := NewStagingArea(client, time.Hour)

	cases := []struct {
		name    string
		setup   func()
		uid     string
		id      string
		wantErr bool
	}{
		{
			name: "existent metadata",
			setup: func() {
				key := fmt.Sprintf("%s/%s", metadata.FakeUid, metadata.FakeMetadata.Title)
				val, err := json.Marshal(metadata.FakeMetadata)
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				err = client.Set(key, val, staging.expiresIn).Err()
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			uid:     metadata.FakeUid,
			id:      metadata.FakeMetadata.Title,
			wantErr: false,
		},
		{
			name:    "non-existent metadata",
			setup:   func() {},
			uid:     metadata.FakeUid,
			id:      metadata.FakeMetadata.Title,
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			t.Cleanup(func() { cleanupStagingArea(client) })

			err := staging.MetadataRemove(tt.uid, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("got nil, want err")
				}
			} else {
				if err != nil {
					t.Fatalf("got err: %v, want nil", err)
				}
			}
		})
	}
}
