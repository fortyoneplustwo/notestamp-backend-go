package auth

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
)

var (
	testToken = "testToken"
	client    *firestore.Client
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

func TestRevokedTokenAdd(t *testing.T) {
	cases := []struct {
		name       string
		setup      func(c *firestore.Client)
		inputToken string
		wantErr    bool
	}{
		{
			name:       "Add revoked token",
			setup:      func(c *firestore.Client) {},
			inputToken: testToken,
			wantErr:    false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() {
				_, _ = client.Collection("revoked").Doc(testToken).Delete(context.TODO())
			})

			store := NewRevokedTokenStore(client, os.Getenv("REVOKED_COLLECTION"))
			err := store.RevokedTokenAdd(testToken, time.Now())

			if tt.wantErr {
				if err == nil {
					t.Fatalf("%s: got nil, want err", tt.name)
				}
			} else {
				if err != nil {
					t.Fatalf("%s: got error %v, want nil", tt.name, err)
				}
			}
		})
	}
}

func TestRevokedTokenCheck(t *testing.T) {
	cases := []struct {
		name       string
		setup      func(c *firestore.Client)
		inputToken string
		wantErr    bool
		wantFound  bool
	}{
		{
			name: "Existent token",
			setup: func(c *firestore.Client) {
				_, _ = client.Collection("revoked").Doc(testToken).Set(context.TODO(), map[string]time.Time{
					"expiry": time.Now(),
				})
			},
			inputToken: testToken,
			wantErr:    false,
			wantFound:  true,
		},
		{
			name:       "Non-existent token",
			setup:      func(c *firestore.Client) {},
			inputToken: testToken,
			wantErr:    false,
			wantFound:  false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() {
				_, _ = client.Collection("revoked").Doc(testToken).Delete(context.TODO())
			})

			store := NewRevokedTokenStore(client, os.Getenv("REVOKED_COLLECTION"))
			found, err := store.RevokedTokenCheck(testToken)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("%s: got nil, want err", tt.name)
				}
			} else {
				if err != nil {
					t.Fatalf("%s: got error %v, want nil", tt.name, err)
				}
				if found != tt.wantFound {
					t.Fatalf("%s: got %v want, %v", tt.name, found, tt.wantFound)
				}
			}
		})
	}
}
