package metadata

type format string

const (
	Audio   format = "audio"
	Youtube format = "youtube"
	Pdf     format = "pdf"
)

var (
	ValidFormats = map[format]string{
		Audio:   "audio/*",
		Youtube: "",
		Pdf:     "application/*",
	}
)

type Metadata struct {
	Format        format `json:"type" firestore:"type" redis:"type"`
	Title         string `json:"title" firestore:"title" redis:"title"`
	Src           string `json:"src" firestore:"src" redis:"src"`
	NotesSize     int    `json:"notesSize" firestore:"notesSize" redis:"notesSize"`
	MediaSize     int    `json:"mediaSize" firestore:"mediaSize" redis:"mediaSize"`
	NotesMimeType string `json:"notesMimeType" firestore:"notesMimeType" redis:"notesMimeType"`
	MediaMimeType string `json:"mediaMimeType" firestore:"mediaMimeType" redis:"mediaMimeType"`
}

func (m Metadata) IsValid() bool {
	if _, ok := ValidFormats[m.Format]; !ok {
		return false
	}
	if m.Title == "" {
		return false
	}
	if m.NotesSize == 0 {
		return false
	}
	if m.NotesMimeType == "" {
		return false
	}
	return true
}
