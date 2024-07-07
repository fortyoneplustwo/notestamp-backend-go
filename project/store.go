package project

type MetadataStore interface {
	Add(uid int, p Metadata) error
	Get(uid int, title string) (Metadata, error)
	List(uid int) ([]string, error)
	Remove(uid int, title string) (Metadata, error)
}

type MediaStore interface {
	Add(uid int, m Media) error
	Get(uid int, title string) (Media, error)
	Stream(uid int, title string) (string, error)
	Remove(uid int, title string) error
}

type NotesStore interface {
	Add(uid int, n Notes) error
	Get(uid int, title string) (Notes, error)
	Remove(uid int, title string) error
}

