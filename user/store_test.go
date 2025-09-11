package user

import (
	"context"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
)

var (
	testUid   = ""
	testEmail = FakeCredentials.Email
	testPwd   = FakeCredentials.Password

	client *firestore.Client
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

func TestUserAdd(t *testing.T) {
	cases := []struct {
		name       string
		setup      func(*firestore.Client)
		inputEmail string
		inputPwd   string
		wantErr    bool
	}{
		{
			name:       "New user",
			setup:      func(c *firestore.Client) {},
			inputEmail: testEmail,
			inputPwd:   testPwd,
			wantErr:    false,
		},
		{
			name: "Existent user",
			setup: func(c *firestore.Client) {
				docRef, _, _ := client.Collection("users").Add(
					context.TODO(),
					Credentials{testEmail, testPwd},
				)
				testUid = docRef.ID
			},
			inputEmail: testEmail,
			inputPwd:   testPwd,
			wantErr:    true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() {
				_, _ = client.Collection("users").Doc(testUid).Delete(context.TODO())
			})

			store := NewUserStore(client, os.Getenv("USER_COLLECTION"))
			uid, err := store.UserAdd(Credentials{tt.inputEmail, tt.inputPwd})

			if tt.wantErr {
				if err == nil {
					t.Fatalf("%s: got nil, want error", tt.name)
				}
			} else {
				if err != nil {
					t.Fatalf(`%s: got %v, want nil`, tt.name, err)
				}
				testUid = uid
			}
		})
	}
}

func TestUserGetByEmail(t *testing.T) {
	cases := []struct {
		name       string
		setup      func(*firestore.Client)
		inputEmail string
		wantErr    bool
	}{
		{
			name:       "Non-existent user",
			setup:      func(c *firestore.Client) {},
			inputEmail: testEmail,
			wantErr:    true,
		},
		{
			name: "Existent user",
			setup: func(c *firestore.Client) {
				docRef, _, _ := client.Collection("users").Add(
					context.TODO(),
					Credentials{testEmail, testPwd},
				)
				testUid = docRef.ID
			},
			inputEmail: testEmail,
			wantErr:    false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() {
				_, _ = client.Collection("users").Doc(testUid).Delete(context.TODO())
			})

			store := NewUserStore(client, os.Getenv("USER_COLLECTION"))
			_, err := store.UserGetByEmail(testEmail)

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

func TestUserGetById(t *testing.T) {
	cases := []struct {
		name     string
		setup    func(*firestore.Client)
		inputUid string
		wantErr  bool
	}{
		{
			name:     "Non-existent user",
			setup:    func(c *firestore.Client) {},
			inputUid: "someUid",
			wantErr:  true,
		},
		{
			name: "Existent user",
			setup: func(c *firestore.Client) {
				docRef, _, _ := client.Collection("users").Add(
					context.TODO(),
					Credentials{testEmail, testPwd},
				)
				testUid = docRef.ID
			},
			inputUid: testUid,
			wantErr:  false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() {
				_, _ = client.Collection("users").Doc(testUid).Delete(context.TODO())
			})

			store := NewUserStore(client, os.Getenv("USER_COLLECTION"))
			user, err := store.UserGetById(testUid)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("%s: got nil, want error", tt.name)
				}
			} else {
				if err != nil {
					t.Fatalf(`%s: got %v, want nil`, tt.name, err)
				}
				testUid = user.Uid
			}
		})
	}
}

func TestUserRemove(t *testing.T) {
	cases := []struct {
		name     string
		setup    func(*firestore.Client)
		inputUid string
		wantErr  bool
	}{
		{
			name:     "Non-existent user",
			setup:    func(c *firestore.Client) {},
			inputUid: "someUid",
			wantErr:  false, // won't return an error even if user does not exist
		},
		{
			name: "Existent user",
			setup: func(c *firestore.Client) {
				docRef, _, _ := client.Collection("users").Add(
					context.TODO(),
					Credentials{testEmail, testPwd},
				)
				testUid = docRef.ID
			},
			inputUid: testUid,
			wantErr:  false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(client)
			t.Cleanup(func() {
				_, _ = client.Collection("users").Doc(testUid).Delete(context.TODO())
			})

			store := NewUserStore(client, os.Getenv("USER_COLLECTION"))
			err := store.UserRemove(testUid)

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
