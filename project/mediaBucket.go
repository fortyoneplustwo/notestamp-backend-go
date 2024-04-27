package project

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MediaBucket struct {
  bucketName string
  client *s3.Client
}


func NewMediaBucket(bucketName string, s *s3.Client) *MediaBucket {
  return &MediaBucket{
    bucketName: bucketName,
    client: s,
  }
}


// Implement MediaStore interface
func (b *MediaBucket) Add(uid int, m Media) error {
  key := strconv.Itoa(uid) + "/" + m.Title

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


func (b *MediaBucket) Get(uid int, title string) (Media, error) {
  key := strconv.Itoa(uid) + "/" + title

  result, err := b.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
	})
  if err != nil {
    return Media{}, err
  }
  defer result.Body.Close()

  return Media{Title: title, Data: result.Body}, nil
}


func (b *MediaBucket) Stream(uid int, title string) (string, error) {
  key := strconv.Itoa(uid) + "/" + title

  presignClient := s3.NewPresignClient(b.client)

  request, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(time.Hour * 60 * 60)
	})
  if err != nil {
    return "", err
  }

  return request.URL, nil
}


func (b *MediaBucket) Remove(uid int, title string) error {
  key := strconv.Itoa(uid) + "/" + title

  _, err := b.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(b.bucketName),
    Key: aws.String(key),
	})
  if err != nil {
    return err
  }

  return nil
}


