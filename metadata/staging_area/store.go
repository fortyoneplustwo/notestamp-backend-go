package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"notestamp/metadata"
	"time"

	"github.com/go-redis/redis"
)

type StagingArea struct {
	client    *redis.Client
	ctx       context.Context
	expiresIn time.Duration
}

func NewStagingArea(c *redis.Client, exp time.Duration) *StagingArea {
	return &StagingArea{
		client:    c,
		ctx:       context.Background(),
		expiresIn: exp,
	}
}

func (s *StagingArea) MetadataAdd(uid string, m metadata.Metadata) error {
	key := fmt.Sprintf("%s/%s", uid, m.Title)
	val, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = s.client.Set(key, val, s.expiresIn).Err()
	return err
}

func (s *StagingArea) MetadataGet(uid string, id string) (metadata.Metadata, error) {
	key := fmt.Sprintf("%s/%s", uid, id)
	val, err := s.client.Get(key).Result()
	if err == redis.Nil {
		return metadata.Metadata{}, metadata.ErrNotExist
	}
	if err != nil {
		return metadata.Metadata{}, err
	}

	var m metadata.Metadata
	err = json.Unmarshal([]byte(val), &m)
	if err != nil {
		return metadata.Metadata{}, err
	}
	return m, nil
}

func (s *StagingArea) MetadataRemove(uid string, id string) error {
	key := fmt.Sprintf("%s/%s", uid, id)
	_, err := s.client.Del(key).Result()
	return err
}

func (s *StagingArea) MetadataRemoveAll(uid string) error {
	return nil
}
