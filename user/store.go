package user

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
)

var ErrUserNotFound error = errors.New("document does not exist")
var ErrUserAlreadyExists error = errors.New("user already exists")

type UserStore struct {
	client *firestore.Client
	ctx    context.Context
}

func NewUserStore(c *firestore.Client) *UserStore {
	return &UserStore{
		client: c,
		ctx:    context.Background(),
	}
}

func (s *UserStore) UserGetById(uid string) (u User, err error) {
	path := fmt.Sprintf("users/%s", uid)
	docRef := s.client.Doc(path)
	snapshot, err := docRef.Get(s.ctx)
	if err != nil {
		return u, err
	}
	err = snapshot.DataTo(&u)
	if err != nil {
		return u, nil
	}
	u.Uid = docRef.ID
	return u, nil
}

func (s *UserStore) UserGetByEmail(email string) (u User, err error) {
	query := s.client.Collection("users").Where("email", "==", email).Limit(1)
	docs, err := query.Documents(s.ctx).GetAll()
	if err != nil {
		return u, err
	}
	if len(docs) == 0 {
		return u, ErrUserNotFound
	}
	err = docs[0].DataTo(&u)
	if err != nil {
		return u, err
	}
	u.Uid = docs[0].Ref.ID
	return u, nil
}

func (s *UserStore) UserAdd(c Credentials) (string, error) {
	_, err := s.UserGetByEmail(c.Email)
	if err == nil {
		return "", ErrUserAlreadyExists
	}
	if !errors.Is(err, ErrUserNotFound) {
		return "", err
	}
	docRef, _, err := s.client.Collection("users").Add(s.ctx, c)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

func (s *UserStore) UserRemove(uid string) error {
	path := fmt.Sprintf("users/%s", uid)
	_, err := s.client.Doc(path).Delete(s.ctx)
	if err != nil {
		return err
	}
	return nil
}
