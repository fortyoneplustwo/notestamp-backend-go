package auth

import (
	"encoding/json"
	"net/http"
	"notestamp/user"
	"runtime"
	"time"
	"github.com/alexedwards/argon2id"
)

type AuthHandler struct{
  userStore user.UserStore
  revokedStore RevokedStore
}

// Constructor
func NewAuthHandler(us user.UserStore, rs RevokedStore) *AuthHandler {
  return &AuthHandler{userStore: us, revokedStore: rs}
}


// Routes
func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  println("inside auth handler")
  w.Write([]byte("This is my auth handler"))
}


func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
  // Validation
  email := r.FormValue("username")
  pw := r.FormValue("password")
  if email == "" || pw == "" {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("Missing form fields"))
    return
  }

  // Generate salt and hash
  params := argon2id.Params{ 
    Iterations: 4,
    Memory: 1024 * 19,
    Parallelism: uint8(runtime.NumCPU()),
    KeyLength: 32,
    SaltLength: 16,
  }
  hash, err := argon2id.CreateHash(pw, &params)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  // Add user to database
  err = h.userStore.Add(user.User{Email: email, Password: hash, Directory: nil})
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }
}


func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
  // Validation
  email := r.FormValue("username")
  pw := r.FormValue("password")
  if email == "" || pw == "" {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("Missing form fields"))
    return
  }

  // Verify password
  user, err := h.userStore.Get(email)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }
  match, err := argon2id.ComparePasswordAndHash(pw, user.Password)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  if !match {
    http.Error(w, "Incorrect password", http.StatusUnauthorized)
    return
  }

  // Generate jwt tokens
  accToken, err := GenerateSignedToken(user.Id, time.Now().Add(time.Hour))
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  refToken, err := GenerateSignedToken(user.Id, time.Now().Add(time.Hour * 24 * 5))
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // Set response cookies
  accCookie := http.Cookie{ Name: "access-token", Value: accToken }
  refCookie := http.Cookie{ Name: "refresh-token", Value: refToken }
  http.SetCookie(w, &accCookie)
  http.SetCookie(w, &refCookie)

  // Set response payload
  payload, err := json.Marshal(user)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(payload)
}


func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
  // Validation
  accToken, err := r.Cookie("access-token")
  if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
  }
  refToken, err := r.Cookie("refresh-token")
  if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
  }

  // Verify access token
  id, exp, err := VerifyToken(accToken.Value)
  if err != nil {
    if err == ErrTokenExpired {
      http.Error(w, err.Error(), http.StatusUnauthorized)
    } else {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    return
  }

  // Revoke refresh token
  if err := h.revokedStore.Add(refToken.Value, id, exp); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  // Cleanup expired tokens (concurrently) from revoked db 
  go h.revokedStore.Cleanup(id)

  // Set new cookies to empty & expired tokens
  accCookie := http.Cookie{ 
    Name: "access-token", Value: "", Expires: time.Unix(0, 0),
  }
  refCookie := http.Cookie{ 
    Name: "refresh-token", Value: "", Expires: time.Unix(0, 0),
  }
  http.SetCookie(w, &accCookie)
  http.SetCookie(w, &refCookie)
}


func (h *AuthHandler) Unregister(w http.ResponseWriter, r *http.Request) {
}

