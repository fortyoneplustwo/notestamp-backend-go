package project

import (
	"bytes"
	"encoding/json"

	"io"
	"net/http"
	"notestamp/auth"
	"notestamp/user"
	"time"

	"github.com/gorilla/mux"
)

type ProjectHandler struct {
	metadataStore MetadataStore
	userStore    user.UserStore
	mediaStore   MediaStore
	notesStore   NotesStore
	revokedStore auth.RevokedStore
}

// Constructor
func NewProjectHandler(ps MetadataStore,
	us user.UserStore,
	ms MediaStore,
	ns NotesStore,
	rs auth.RevokedStore) *ProjectHandler {
	return &ProjectHandler{
		metadataStore: ps,
		userStore:    us,
		mediaStore:   ms,
		notesStore:   ns,
		revokedStore: rs,
	}
}

// Routes
func (h *ProjectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my project handler"))
}

func (h *ProjectHandler) Save(w http.ResponseWriter, r *http.Request) {
	// Validate request
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

	mediaReceived := true
	mediaFile, _, err := r.FormFile("mediaFile")
	if err != nil {
    if err == http.ErrMissingFile {
      mediaReceived = false
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
  if mediaReceived {
    defer mediaFile.Close()
  }

	notesFile, _, err := r.FormFile("notesFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
  defer notesFile.Close()

	metadata := r.FormValue("metadata")
	if metadata == "" {
		http.Error(w, "Missing metadata", http.StatusBadRequest)
		return
	}

	// Process form data into our custom data types
	project, err := NewMetadata([]byte(metadata))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var media Media
	if mediaReceived {
    project.SetMediaHash(mediaFile)
		media, err = NewMedia(project.MediaHash, project.Mimetype, mediaFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	notes, err := NewNotes(project.Title, notesFile)
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

	// Upload notes to notesStore
	if err := h.notesStore.Add(uid, notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Upload media to mediaStore if no duplicate exists
	if mediaReceived {
    if isDup, _ := h.metadataStore.FindMediaDup(uid, project.MediaHash); !isDup {
      if err := h.mediaStore.Add(uid, media); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
      }
    }
	}

	// Save project to projectStore
	err = h.metadataStore.Add(uid, project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated directory
	dir, err := h.metadataStore.List(uid)
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

	projects, err := p.metadataStore.List(uid)
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
	w.Write(jpayload)
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

	project, err := p.metadataStore.Get(uid, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set media source if not already set
	domain := "http://localhost:8080"
	if project.Src == "" {
		if project.Format == "audio" {
			project.Src = domain + "/media/stream/" + project.MediaHash
		} else {
			project.Src = domain + "/media/download/" + project.MediaHash
		}
	}

	// Get notes file
	notes, err := p.notesStore.Get(uid, project.MakeNotesKey())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

  notesBuff := new(bytes.Buffer)
  if _, err = io.Copy(notesBuff, notes.Data); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  defer notes.Data.Close()

  payload := project.ToMap()
	payload["notes"] = notesBuff.String()

	jpayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jpayload)
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

	deletedProject, err := p.metadataStore.Remove(uid, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = p.notesStore.Remove(uid, deletedProject.MakeNotesKey())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

  if deletedProject.MediaHash != "" {
    err = p.mediaStore.Remove(uid, deletedProject.MakeMediaKey())
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }

	// Return updated directory
	dir, err := p.metadataStore.List(uid)
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

func (h *ProjectHandler) DownloadMedia(w http.ResponseWriter, r *http.Request) {
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
			newToken, err := auth.Refresh(refCookie.Value, time.Now().Add(time.Hour), h.revokedStore)
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

	project, err := h.metadataStore.Get(uid, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	media, err := h.mediaStore.Get(uid, project.MakeMediaKey())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", project.Mimetype)
  _, err = io.Copy(w, media.Data)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  defer media.Data.Close()
}

func (h *ProjectHandler) StreamMedia(w http.ResponseWriter, r *http.Request) {

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
			newToken, err := auth.Refresh(refCookie.Value, time.Now().Add(time.Hour), h.revokedStore)
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

	project, err := h.metadataStore.Get(uid, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url, err := h.mediaStore.Stream(uid, project.MakeMediaKey())

	payload := map[string]string{
		"streamURL": url,
	}
	jpayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jpayload)
}
