package project

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type NotesBucket struct {
	bucketName string
	client     *s3.Client
}

func NewNotesBucket(bucketName string, s *s3.Client) *NotesBucket {
	return &NotesBucket{
		bucketName: bucketName,
		client:     s,
	}
}

// Implement NotesStore interface
func (b *NotesBucket) Add(uid int, m Notes) error {
	key := strconv.Itoa(uid) + "/" + m.Title

	fmt.Println(context.TODO())

	_, err := b.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
		Body:   m.Data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (b *NotesBucket) Get(uid int, title string) (Notes, error) {
	key := strconv.Itoa(uid) + "/" + title

	result, err := b.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return Notes{}, err
	}

	return Notes{Title: title, Data: result.Body}, nil
}

func (b *NotesBucket) Remove(uid int, title string) error {
	key := strconv.Itoa(uid) + "/" + title

	_, err := b.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}
