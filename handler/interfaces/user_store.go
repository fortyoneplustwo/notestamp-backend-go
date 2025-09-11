package handler

import "notestamp/user"


type UserAdder interface {
	UserAdd(u user.Credentials) (string, error)
}

type UserGetter interface {
	UserGetById(uid string) (user.User, error)
	UserGetByEmail(email string) (user.User, error)
}

type UserRemover interface {
	UserRemove(uid string) error
}
