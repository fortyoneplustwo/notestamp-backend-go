package handler

import "time"

type RevokedTokenAdder interface {
	RevokedTokenAdd(token string, exp time.Time) error
}

type RevokedTokenChecker interface {
	RevokedTokenCheck(token string) (bool, error)
}
