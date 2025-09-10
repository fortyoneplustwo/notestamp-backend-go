package handler

import (
	"io"
)

type NotesRemover interface {
	NotesRemove(uid string, oid string) error
	NotesRemoveAll(uid string) error
}

type NotesUpdater interface {
	NotesUpdate(uid string, oid string, f io.Reader) error
}

type NotesUploadPathGetter interface {
	NotesGetUploadPath(uid string, oid string) string
}

type NotesChecker interface {
	NotesCheck(uid string, oid string) (bool, error)
}

type NotesCheckRemover interface {
	NotesCheck(uid string, oid string) (bool, error)
	NotesRemove(uid string, oid string) error
	NotesRemoveAll(uid string) error
}
