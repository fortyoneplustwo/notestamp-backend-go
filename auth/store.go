package auth

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RevokedTokenStore struct {
	client *firestore.Client
	ctx    context.Context
	collection string
}

func NewRevokedTokenStore(c *firestore.Client, collection string) *RevokedTokenStore {
	return &RevokedTokenStore{
		client: c,
		ctx:    context.Background(),
		collection: collection,
	}
}

func (s RevokedTokenStore) RevokedTokenAdd(token string, exp time.Time) error {
	path := fmt.Sprintf("%s/%s", s.collection, token)
	_, err := s.client.Doc(path).Set(s.ctx, map[string]any{
		"expiry": exp,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s RevokedTokenStore) RevokedTokenCheck(token string) (bool, error) {
	path := fmt.Sprintf("%s/%s", s.collection, token)
	_, err := s.client.Doc(path).Get(s.ctx)
	if status.Code(err) == codes.NotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
