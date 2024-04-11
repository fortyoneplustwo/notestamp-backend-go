package project

import (
	"encoding/json"
	"net/http"
	"notestamp/auth"
	"notestamp/user"
	"time"

	"github.com/gorilla/mux"
)

type ProjectHandler struct {
  projectStore ProjectStore
  userStore user.UserStore
  revokedStore auth.RevokedStore
}


// Constructor
func NewProjectHandler(ps ProjectStore, us user.UserStore, rs auth.RevokedStore) *ProjectHandler {
  return &ProjectHandler {
    projectStore: ps,
    userStore: us,
    revokedStore: rs,
  }
}


// Routes
func (h *ProjectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("This is my project handler"))
}


func (h *ProjectHandler) Save(w http.ResponseWriter, r *http.Request) {
  // Validation
  // _ := r.FormValue("mediaFile")
  notes := r.FormValue("content")
  metadata := r.FormValue("metadata")
  if notes == "" || metadata == "" {
    http.Error(w, "Missing form fields", http.StatusBadRequest)
    return
  }

  accCookie, err := r.Cookie("access-token")
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  refCookie, err := r.Cookie("refresh-token")
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  project := Project{}
  err = json.Unmarshal([]byte(metadata), &project)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  // Authentication
  uid, _, err := auth.VerifyToken(accCookie.Value)
  if err != nil {
    if err != auth.ErrTokenExpired {
      http.Error(w, err.Error(), http.StatusInternalServerError) 
      return
    } else {
      // Access token has expired. Needs a refresh.
      newAccToken, err := auth.Refresh(refCookie.Value, time.Now().Add(time.Hour), h.revokedStore)
      if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
      }
      accCookie.Value = newAccToken
    }
  }

  http.SetCookie(w, accCookie)
  http.SetCookie(w, refCookie)

  // Upload notes to s3

  // Upload media to s3

  // Save project to database
  err = h.projectStore.Add(uid, project)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // Return updated directory
  dir, err := h.projectStore.List(uid)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  payload := map[string][]string{
    "directory": dir,
  }
  jpayload, err := json.Marshal(payload)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(jpayload)
}


func (p *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
  accCookie, err := r.Cookie("access-token")
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  refCookie, err := r.Cookie("refresh-token")
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  uid, _, err := auth.VerifyToken(accCookie.Value)
  if err != nil {
    if err == auth.ErrTokenExpired {
      newToken, err := auth.Refresh(refCookie.Value, time.Now().Add(time.Hour), p.revokedStore)
      if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
      }
      accCookie.Value = newToken
    } else {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }

  http.SetCookie(w, accCookie)
  http.SetCookie(w, refCookie)

  projects, err := p.projectStore.List(uid)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  
  payload := map[string][]string{
    "projects": projects,
  }
  jpayload, err := json.Marshal(payload)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  _, err = w.Write(jpayload)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}


func (p *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
  title := mux.Vars(r)["title"]
  if title == "" {
    http.Error(w, "No title provided", http.StatusBadRequest)
  }

  accCookie, err := r.Cookie("access-token")
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  refCookie, err := r.Cookie("refresh-token")
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  uid, _, err := auth.VerifyToken(accCookie.Value)
  if err != nil {
    if err == auth.ErrTokenExpired {
      newToken, err := auth.Refresh(refCookie.Value, time.Now().Add(time.Hour), p.revokedStore)
      if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
      }
      accCookie.Value = newToken
    } else {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }

  http.SetCookie(w, accCookie)
  http.SetCookie(w, refCookie)

  project, err := p.projectStore.Get(uid, title)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  payload, err := json.Marshal(project)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(payload)
}


func (p *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
  title := mux.Vars(r)["title"]
  if title == "" {
    http.Error(w, "No title provided", http.StatusBadRequest)
  }

  accCookie, err := r.Cookie("access-token")
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  refCookie, err := r.Cookie("refresh-token")
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  uid, _, err := auth.VerifyToken(accCookie.Value)
  if err != nil {
    if err == auth.ErrTokenExpired {
      newToken, err := auth.Refresh(refCookie.Value, time.Now().Add(time.Hour), p.revokedStore)
      if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
      }
      accCookie.Value = newToken
    } else {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }

  http.SetCookie(w, accCookie)
  http.SetCookie(w, refCookie)

  err = p.projectStore.Remove(uid, title)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // Return updated directory
  dir, err := p.projectStore.List(uid)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  payload := map[string][]string{
    "directory": dir,
  }
  jpayload, err := json.Marshal(payload)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(jpayload)
}
