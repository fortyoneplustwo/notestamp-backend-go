package handler

import (
	"encoding/json"
	"net/http"
	handler "notestamp/handler/interfaces"
	"notestamp/metadata"
	"notestamp/user"
)

func ListMetadata(s handler.MetadataLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser, ok := user.FromContext(r.Context())
		if !ok {
			http.Error(w, "could not extract user id", http.StatusUnauthorized)
			return
		}

		list, err := s.MetadataList(authUser.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		payload, err := json.Marshal(
			struct {
				Projects []metadata.Metadata `json:"projects"`
			}{
				Projects: list,
			})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(payload)
	}
}

