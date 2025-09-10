package media

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type MediaStore struct {
	client *storage.Client
	bucket string
	ctx    context.Context
}

func NewMediaStore(c *storage.Client, bucket string) *MediaStore {
	return &MediaStore{
		client: c,
		bucket: bucket,
		ctx:    context.Background(),
	}
}

func (s *MediaStore) MediaGetUploadPath(uid string, oid string) string {
	return fmt.Sprintf("%s/%s/%s/%s", s.bucket, uid, Prefix, oid)
}

func (s *MediaStore) MediaCheck(uid string, oid string) (bool, error) {
	path := fmt.Sprintf("%s/%s/%s", uid, Prefix, oid)
	object := s.client.Bucket(s.bucket).Object(path)
	_, err := object.Attrs(s.ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *MediaStore) MediaRemove(uid string, oid string) error {
	path := fmt.Sprintf("%s/%s/%s", uid, Prefix, oid)
	object := s.client.Bucket(s.bucket).Object(path)
	err := object.Delete(s.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *MediaStore) MediaRemoveAll(uid string) error {
	it := s.client.Bucket(s.bucket).Objects(s.ctx, &storage.Query{
		Prefix:    uid + "/" + Prefix,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Bucket(%q).Objects(): %w", s.bucket, err)
		}
		err = s.client.Bucket(s.bucket).Object(attrs.Name).Delete(s.ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
