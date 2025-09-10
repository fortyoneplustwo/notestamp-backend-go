package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"notestamp/cleanup"
	handler "notestamp/handler/interfaces"
	project "notestamp/metadata"
	"notestamp/notes"
	"notestamp/user"
	"regexp"
)

func DeleteProject(dr handler.MetadataRemover, mr handler.MediaRemover, nr handler.NotesRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser, ok := user.FromContext(r.Context())
		if !ok {
			http.Error(w, "could not authorize user", http.StatusUnauthorized)
			return
		}

		var project struct{ id string }
		err := json.NewDecoder(r.Body).Decode(&project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		targets := []struct{ remove func() error }{
			{remove: func() error { return mr.MediaRemove(authUser.Id, project.id) }},
			{remove: func() error { return nr.NotesRemove(authUser.Id, project.id) }},
		}
		ch := make(chan error, 2)
		for _, target := range targets {
			go func() { ch <- target.remove() }()
		}
		for range targets {
			err := <-ch
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		err = dr.MetadataRemove(authUser.Id, project.id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func StageWithoutMedia(staging handler.MetadataAdder, ns handler.NotesUploadPathGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser, ok := user.FromContext(r.Context())
		if !ok {
			http.Error(w, "failed to extract uid", http.StatusUnauthorized)
			return
		}

		var metadata project.Metadata
		err := json.NewDecoder(r.Body).Decode(&metadata)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if ok := metadata.IsValid(); !ok {
			http.Error(w, "invalid metadata", http.StatusBadRequest)
			return
		}
		if metadata.Src == "" {
			http.Error(w, "missing src", http.StatusBadRequest)
			return
		}
		metadata.MediaMimeType = ""
		metadata.MediaSize = 0

		err = staging.MetadataAdd(authUser.Id, metadata)
		if err != nil && err != project.ErrDuplicate {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		path := ns.NotesGetUploadPath(authUser.Id, metadata.Title)

		w.Header().Add("Content-Type", "application/json")
		payload, err := json.Marshal(struct {
			Key string `json:"key"`
		}{
			Key: path,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(payload)
	}
}

func StageWithMedia(staging handler.MetadataAdder, ms handler.MediaUploadPathGetter, ns handler.NotesUploadPathGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser, ok := user.FromContext(r.Context())
		if !ok {
			http.Error(w, "failed to extract uid", http.StatusUnauthorized)
			return
		}

		var metadata project.Metadata
		err := json.NewDecoder(r.Body).Decode(&metadata)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if ok := metadata.IsValid(); !ok {
			http.Error(w, "invalid metadata", http.StatusBadRequest)
			return
		}
		if metadata.MediaMimeType == "" {
			http.Error(w, "missing mediaMimeType", http.StatusBadRequest)
			return
		}
		if metadata.MediaSize == 0 {
			http.Error(w, "missing mediaSize", http.StatusBadRequest)
			return
		}
		validMediaMimeType := project.ValidFormats[metadata.Format]
		re, err := regexp.Compile(validMediaMimeType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if ok = re.MatchString(metadata.MediaMimeType); !ok {
			http.Error(w, "invalid media MIME type", http.StatusBadRequest)
			return
		}
		metadata.Src = ""

		err = staging.MetadataAdd(authUser.Id, metadata)
		if err != nil && err != project.ErrDuplicate {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		notesPath := ns.NotesGetUploadPath(authUser.Id, metadata.Title)
		mediaPath := ms.MediaGetUploadPath(authUser.Id, metadata.Title)

		w.Header().Add("Content-Type", "application/json")
		payload, err := json.Marshal(struct {
			MediaKey string `json:"mediaKey"`
			NotesKey string `json:"notesKey"`
		}{
			MediaKey: mediaPath,
			NotesKey: notesPath,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(payload)
	}
}

func UpdateProject(nu handler.NotesUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser, ok := user.FromContext(r.Context())
		if !ok {
			http.Error(w, "failed to extract uid", http.StatusUnauthorized)
			return
		}

		var newNotes notes.Notes
		err := json.NewDecoder(r.Body).Decode(&newNotes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err = nu.NotesUpdate(authUser.Id, newNotes.Id, bytes.NewBuffer(newNotes.Data))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func PostSaveProject(staging handler.MetadataGetRemover, ma handler.MetadataAdder, mc handler.MediaCheckRemover, nc handler.NotesCheckRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser, ok := user.FromContext(r.Context())
		if !ok {
			http.Error(w, "failed to extract uid", http.StatusUnauthorized)
			return
		}

		var toCommit struct {
			Id string `json:"id"`
		}
		err := json.NewDecoder(r.Body).Decode(&toCommit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		staged, err := staging.MetadataGet(authUser.Id, toCommit.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		targets := []struct {
			name      string
			check func() (bool, error)
		}{
			{
				name:      "notes",
				check: func() (bool, error) { return nc.NotesCheck(authUser.Id, staged.Title) },
			},
			{
				name:      "media",
				check: func() (bool, error) { return mc.MediaCheck(authUser.Id, staged.Title) },
			},
		}

		type result struct {
			ok  bool
			err error
		}
		ch := make(chan result, 2)

		for _, target := range targets {
			if target.name == "media" && staged.Src != "" {
				continue
			}
			go func() {
				ok, err := target.check()
				ch <- result{ok, err}
			}()
		}

		for _, t := range targets {
			if t.name == "media" && staged.Src != "" {
				continue
			}
			result := <-ch
			if result.err != nil {
				http.Error(w, result.err.Error(), http.StatusInternalServerError)
				return
			}
			if !result.ok {
				hasMedia := staged.Src == ""
				go cleanup.FailedSave(staging, mc, nc, authUser.Id, staged.Title, &cleanup.Options{HasMedia: hasMedia})
				http.Error(w, "failed to save", http.StatusInternalServerError)
				return
			}
		}

		err = ma.MetadataAdd(authUser.Id, staged)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = staging.MetadataRemove(authUser.Id, staged.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
