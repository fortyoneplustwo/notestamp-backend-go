package media

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	fakeOid   = "fakeOid"
	fakeUid   = "fakeUid"
	fakeMedia = bytes.NewBuffer([]byte("fakeMedia"))

	client *storage.Client
)

func clearBucket(client *storage.Client, bucket string, uid string) {
	it := client.Bucket(bucket).Objects(context.TODO(), &storage.Query{
		Prefix: uid + "/" + Prefix,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		_ = client.Bucket(bucket).Object(attrs.Name).Delete(context.TODO())
	}
}

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	path := os.Getenv("FIREBASE_CONF")
	ctx := context.Background()
	c, err := storage.NewClient(ctx, option.WithCredentialsFile(path))
	if err != nil {
		log.Fatalf("failed to init client: %v", err)
	}
	client = c

	exitCode := m.Run()

	err = c.Close()
	if err != nil {
		log.Printf("failed to close client: %v", err)
	}

	os.Exit(exitCode)
}

func TestMediaCheck(t *testing.T) {
	bucketId := os.Getenv("MEDIA_BUCKET")

	cases := []struct {
		name      string
		setup     func(*storage.Client)
		uid       string
		oid       string
		wantError bool
		wantFound bool
	}{
		{
			name: "existent media",
			setup: func(c *storage.Client) {
				path := fmt.Sprintf("%s/%s/%s", fakeUid, Prefix, fakeOid)
				object := c.Bucket(bucketId).Object(path).NewWriter(context.TODO())
				if _, err := io.Copy(object, fakeMedia); err != nil {
					object.Close()
					t.Fatalf("failed to setup test: %v", err)
				}
				if err := object.Close(); err != nil {
					t.Fatalf("failed to setup test: %v", err)
				}
			},
			uid:       fakeUid,
			oid:       fakeOid,
			wantError: false,
			wantFound: true,
		},
		{
			name:      "non-existent media",
			setup:     func(c *storage.Client) {},
			uid:       fakeUid,
			oid:       fakeOid,
			wantError: false,
			wantFound: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() { clearBucket(client, bucketId, tt.uid) })

			mediaStore := NewMediaStore(client, bucketId)
			found, err := mediaStore.MediaCheck(tt.uid, tt.oid)

			if tt.wantError {
				if err == nil {
					t.Fatalf("got nil, want err")
				}
			} else {
				if found != tt.wantFound {
					t.Fatalf("got %v, want %v", found, tt.wantFound)
				}
			}
		})
	}
}

func TestMediaRemove(t *testing.T) {
	bucketId := os.Getenv("MEDIA_BUCKET")

	cases := []struct {
		name    string
		setup   func(c *storage.Client)
		uid     string
		oid     string
		wantErr bool
	}{
		{
			name: "existent media",
			setup: func(c *storage.Client) {
				path := fmt.Sprintf("%s/%s/%s", fakeUid, Prefix, fakeOid)
				object := c.Bucket(bucketId).Object(path).NewWriter(context.TODO())
				if _, err := io.Copy(object, fakeMedia); err != nil {
					object.Close()
					t.Fatalf("failed to setup test: %v", err)
				}
				if err := object.Close(); err != nil {
					t.Fatalf("failed to setup test: %v", err)
				}
			},
			uid:     fakeUid,
			oid:     fakeOid,
			wantErr: false,
		},
		{
			name:    "non-existent media",
			setup:   func(c *storage.Client) {},
			uid:     fakeUid,
			oid:     fakeOid,
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() { clearBucket(client, bucketId, tt.uid) })

			mediaStore := NewMediaStore(client, bucketId)
			err := mediaStore.MediaRemove(tt.uid, tt.oid)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("got nil, want err")
				}
			} else {
				if err != nil {
					t.Fatalf("got %v, want nil", err)
				}
			}
		})
	}
}

func TestMediaRemoveAll(t *testing.T) {
	bucketId := os.Getenv("MEDIA_BUCKET")

	cases := []struct {
		name    string
		setup   func(c *storage.Client)
		uid     string
		oid     string
		wantErr bool
	}{
		{
			name: "2 files",
			setup: func(c *storage.Client) {
				for i := range 2 {
					count := strconv.Itoa(i)
					path := fmt.Sprintf("%s/%s/%s", fakeUid, Prefix, fakeOid+count)
					object := c.Bucket(bucketId).Object(path).NewWriter(context.TODO())
					if _, err := io.Copy(object, fakeMedia); err != nil {
						object.Close()
						t.Fatalf("failed to setup test: %v", err)
					}
					if err := object.Close(); err != nil {
						t.Fatalf("failed to setup test: %v", err)
					}
				}
			},
			uid:     fakeUid,
			oid:     fakeOid,
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() { clearBucket(client, bucketId, tt.uid) })

			mediaStore := NewMediaStore(client, bucketId)
			err := mediaStore.MediaRemoveAll(tt.uid)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("got nil, want err")
				}
			} else {
				if err != nil {
					t.Fatalf("got %v, want nil", err)
				}
			}
		})
	}

}
