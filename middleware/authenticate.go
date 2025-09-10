package middleware

import (
	"net/http"
	jwt "notestamp/auth/jwt"
	"notestamp/user"
)

type RevokedTokenChecker interface {
	RevokedTokenCheck(token string) (bool, error)
}

func Authenticate(h http.Handler, rtc RevokedTokenChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acc, ref, err := jwt.ExtractTokens(r)
		if err == http.ErrNoCookie {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		uid, _, err := jwt.VerifyToken(acc)

		switch {
		case err != jwt.ErrTokenExpired:
			accExpTime, refExpTime, err := jwt.CalcTokenExpTime("../.env")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			isRevoked, err := rtc.RevokedTokenCheck(ref)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if isRevoked {
				http.Error(w, "token is revoked already", http.StatusUnauthorized)
				return
			}
			acc, ref, err = jwt.Refresh(ref, accExpTime, refExpTime)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
		case err != nil:
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := user.NewContext(r.Context(), &user.AuthUser{Id: uid})
		jwt.SetTokensAsCookies(w, acc, ref)

		h.ServeHTTP(w, r.WithContext(ctx))
	}
}
