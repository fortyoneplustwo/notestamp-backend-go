package auth

import "time"

type MockRevokedStore struct {
	Err       error
	IsRevoked bool
}

func (m MockRevokedStore) RevokedTokenAdd(token string, exp time.Time) error {
	return m.Err
}

func (m MockRevokedStore) RevokedTokenCheck(token string) (bool, error) {
	return m.IsRevoked, m.Err
}
