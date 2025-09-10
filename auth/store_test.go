package auth

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
)

var (
	testToken = "testToken"
)

func initClient() (*firestore.Client, error) {
	if err := godotenv.Load("../.env"); err != nil {
		return nil, err
	}
	ctx := context.TODO()
	client, err := firestore.NewClient(ctx, os.Getenv("FIREBASE_PROJECT_ID"))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func closeClient(c *firestore.Client) {
	if err := c.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error closing client: %q", err)
	}
}

func TestRevokedTokenAdd(t *testing.T) {
	client, err := initClient()
	if err != nil {
		t.Error(err)
	}
	defer closeClient(client)
	store := NewRevokedTokenStore(client)

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

			err = store.RevokedTokenAdd(testToken, time.Now())

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
	client, err := initClient()
	if err != nil {
		t.Error(err)
	}
	defer closeClient(client)
	store := NewRevokedTokenStore(client)

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
