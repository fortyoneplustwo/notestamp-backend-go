package project

type MediaStore interface {
	Add(uid int, m Media) error
	Get(uid int, title string) (Media, error)
	Stream(uid int, title string) (string, error)
	Remove(uid int, title string) error
}
