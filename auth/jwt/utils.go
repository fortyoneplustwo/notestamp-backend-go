package auth

import "net/http"

func ExtractTokens(r *http.Request) (acc string, ref string, err error) {
	accCookie, err := r.Cookie("access-token")
	if err != nil {
		return acc, ref, err
	}
	refCookie, err := r.Cookie("refresh-token")
	if err != nil {
		return acc, ref, err
	}
	return accCookie.Value, refCookie.Value, nil
}

func SetTokensAsCookies(w http.ResponseWriter, a string, r string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		Value:    a,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    r,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func InvalidateTokenCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
