package metadata

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MetadataStore struct {
	client     *firestore.Client
	ctx        context.Context
	collection string
}

func NewMetadataStore(c *firestore.Client, collection string) *MetadataStore {
	return &MetadataStore{
		client:     c,
		ctx:        context.Background(),
		collection: collection,
	}
}

func (s *MetadataStore) MetadataGet(uid string, id string) (
	m Metadata,
	err error,
) {
	path := fmt.Sprintf("%s/%s", s.collection, uid)
	docRef := s.client.Doc(path)
	snapshot, err := docRef.Get(s.ctx)
	if err != nil {
		return m, err
	}
	metadataMap := make(map[string]Metadata)
	err = snapshot.DataTo(&metadataMap)
	if err != nil {
		return m, err
	}
	return metadataMap[id], nil
}

func (s *MetadataStore) MetadataList(uid string) (l []Metadata, err error) {
	path := fmt.Sprintf("%s/%s", s.collection, uid)
	docRef := s.client.Doc(path)
	snapshot, err := docRef.Get(s.ctx)
	if status.Code(err) == codes.NotFound {
		return l, nil
	}
	if err != nil {
		return l, err
	}
	metadataMap := make(map[string]Metadata)
	err = snapshot.DataTo(&metadataMap)
	if err != nil {
		return l, err
	}
	for _, val := range metadataMap {
		l = append(l, val)
	}
	return l, nil
}

func (s *MetadataStore) MetadataAdd(uid string, m Metadata) error {
	path := fmt.Sprintf("%s/%s", s.collection, uid)
	docRef := s.client.Doc(path)
	snapshot, err := docRef.Get(s.ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		return err
	}
	if status.Code(err) == codes.NotFound {
		_, err = s.client.Doc(path).Set(s.ctx, map[string]Metadata{m.Title: m})
		if err != nil {
			return err
		}
		return nil
	}

	metadataMap := make(map[string]Metadata)
	err = snapshot.DataTo(&metadataMap)
	if err != nil {
		return err
	}
	if _, ok := metadataMap[m.Title]; ok {
		return ErrDuplicate
	}

	_, err = s.client.Doc(path).Set(s.ctx, map[string]Metadata{
		m.Title: m,
	}, firestore.MergeAll)
	if err != nil {
		return err
	}
	return nil
}

func (s *MetadataStore) MetadataUpdate(uid string, m Metadata) error {
	path := fmt.Sprintf("%s/%s", s.collection, uid)
	_, err := s.client.Doc(path).Update(s.ctx, []firestore.Update{
		{Path: m.Title, Value: m},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *MetadataStore) MetadataRemove(uid string, id string) error {
	path := fmt.Sprintf("%s/%s", s.collection, uid)
	_, err := s.client.Doc(path).Update(s.ctx, []firestore.Update{
		{Path: id, Value: firestore.Delete},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *MetadataStore) MetadataRemoveAll(uid string) error {
	_, err := s.client.Collection(s.collection).Doc(uid).Delete(s.ctx)
	if err != nil {
		return err
	}
	return nil
}
