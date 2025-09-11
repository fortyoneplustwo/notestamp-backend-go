package handler

import (
	"encoding/json"
	"net/http"
	"notestamp/auth"
	jwt "notestamp/auth/jwt"
	store "notestamp/handler/interfaces"
	"notestamp/user"

	"github.com/alexedwards/argon2id"
)

func RegisterUser(s store.UserAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		credentials := user.Credentials{}
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		hash, err := auth.CreateHash(credentials.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = s.UserAdd(user.Credentials{
			Email:    credentials.Email,
			Password: hash,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func DeregisterUser(ur store.UserRemover, mr store.MetadataRemover, dr store.MediaRemover, nr store.NotesRemover, ra store.RevokedTokenAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accToken, refToken, err := jwt.ExtractTokens(r)
		if err == http.ErrNoCookie {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		uid, _, err := jwt.VerifyToken(accToken)
		switch {
		case err == jwt.ErrTokenExpired:
			_, _, err = jwt.VerifyToken(refToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
		case err != nil:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ch := make(chan error, 3)
		targets := []struct{ removeAll func() error }{
			{removeAll: func() error { return mr.MetadataRemoveAll((uid)) }},
			{removeAll: func() error { return dr.MediaRemoveAll((uid)) }},
			{removeAll: func() error { return nr.NotesRemoveAll(uid) }},
		}
		for _, target := range targets {
			go func() { ch <- target.removeAll() }()
		}
		for range targets {
			err := <-ch
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		err = ur.UserRemove(uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, refExpTime, err := jwt.ParseTokenWithoutValidation(refToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = ra.RevokedTokenAdd(refToken, refExpTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jwt.InvalidateTokenCookies(w)

		w.WriteHeader(http.StatusNoContent)
	}
}

func LoginUser(ug store.UserGetter, ra store.RevokedTokenAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		credentials := user.Credentials{}
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var refTokenSentWithReq bool
		var oldRefToken string
		refCookie, err := r.Cookie("refresh-token")
		if err == http.ErrNoCookie {
			oldRefToken = ""
			refTokenSentWithReq = false
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			oldRefToken = refCookie.Value
			refTokenSentWithReq = true
		}

		user, err := ug.UserGetByEmail(credentials.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		match, err := argon2id.ComparePasswordAndHash(
			credentials.Password,
			user.Password,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !match {
			http.Error(w, credentials.Password, http.StatusUnauthorized)
			return
		}

		accExpTime, refExpTime, err := jwt.CalcTokenExpTime("../.env")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		accToken, refToken, err := jwt.GenerateTokens(user.Uid, accExpTime, refExpTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if refTokenSentWithReq {
			err = ra.RevokedTokenAdd(oldRefToken, refExpTime)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		jwt.SetTokensAsCookies(w, accToken, refToken)

		payload, err := json.Marshal(struct {
			Email string `json:"email"`
		}{
			Email: user.Email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write(payload)
	}
}

func LogoutUser(ra store.RevokedTokenAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accToken, refToken, err := jwt.ExtractTokens(r)
		if err == http.ErrNoCookie {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _, err = jwt.VerifyToken(accToken)
		switch {
		case err == jwt.ErrTokenExpired:
			_, _, err = jwt.VerifyToken(refToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
		case err != nil:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, refExpTime, err := jwt.ParseTokenWithoutValidation(refToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = ra.RevokedTokenAdd(refToken, refExpTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		jwt.InvalidateTokenCookies(w)

		w.WriteHeader(http.StatusOK)
	}
}
