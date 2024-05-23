package project

type NotesStore interface {
	Add(uid int, n Notes) error
	Get(uid int, title string) (Notes, error)
	Remove(uid int, title string) error
}
