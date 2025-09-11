package auth

import (
	"runtime"

	"github.com/alexedwards/argon2id"
)

var (
	params = argon2id.Params{
		Iterations:  4,
		Memory:      1024 * 19,
		Parallelism: uint8(runtime.NumCPU()),
		KeyLength:   32,
		SaltLength:  16,
	}
)

func CreateHash(pwd string) (hash string, err error) {
	hash, err = argon2id.CreateHash(pwd, &params)
	if err != nil {
		return hash, err
	}
	return hash, nil
}
