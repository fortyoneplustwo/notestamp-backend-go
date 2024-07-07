package auth

import (
  "time"
)

type RevokedStore interface {
	Add(string, int, time.Time) error
	IsRevoked(string, int) (bool, error)
	Cleanup(int)
}
