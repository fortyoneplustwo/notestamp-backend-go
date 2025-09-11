package user

import "context"

type contextKey string

type AuthUser struct {
	Id string
}

var (
	authUserKey = contextKey("authUser")
)

func NewContext(ctx context.Context, authUser *AuthUser) context.Context {
	return context.WithValue(ctx, authUserKey, authUser)
}

func FromContext(ctx context.Context) (AuthUser, bool) {
	authUser, ok := ctx.Value(authUserKey).(*AuthUser)
	if !ok || authUser == nil {
		return AuthUser{}, false
	}
	return *authUser, ok
}
