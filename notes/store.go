package notes

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type NotesStore struct {
	client *storage.Client
	bucket string
	ctx    context.Context
}

func NewNotesStore(c *storage.Client, bucket string) *NotesStore {
	return &NotesStore{
		client: c,
		bucket: bucket,
		ctx:    context.Background(),
	}
}

func (s *NotesStore) NotesGetUploadPath(uid string, oid string) string {
	 return fmt.Sprintf("%s/%s/%s/%s", s.bucket, uid, Prefix, oid)
}

func (s *NotesStore) NotesCheck(uid string, oid string) (bool, error) {
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

func (s *NotesStore) NotesUpdate(uid string, oid string, f io.Reader) error {
	path := fmt.Sprintf("%s/%s/%s", uid, Prefix, oid)
	object := s.client.Bucket(s.bucket).Object(path)
	_, err := object.Attrs(s.ctx)
	if err == storage.ErrObjectNotExist {
		return err
	}
	if err != nil {
		return err
	}
	w := object.NewWriter(s.ctx)
	w.ContentType = ContentType
	if _, err := io.Copy(w, f); err != nil {
		w.Close()
		return err
	}
	if err := w.Close(); err != nil {
		return nil
	}
	return nil
}

func (s *NotesStore) NotesRemove(uid string, oid string) error {
	path := fmt.Sprintf("%s/%s/%s", uid, Prefix, oid)
	object := s.client.Bucket(s.bucket).Object(path)
	err := object.Delete(s.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *NotesStore) NotesRemoveAll(uid string) error {
	it := s.client.Bucket(s.bucket).Objects(s.ctx, &storage.Query{
		Prefix: uid + "/" + Prefix,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		err = s.client.Bucket(s.bucket).Object(attrs.Name).Delete(s.ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
