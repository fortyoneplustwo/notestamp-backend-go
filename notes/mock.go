package notes

import "io"

type MockNotesStore struct {
	Err error
}

var (
	FakeNotes      = "fakeNotes"
	FakeUploadPath = "fakeUploadPath"
)

func (m MockNotesStore) NotesRemove(uid string, oid string) error {
	return m.Err
}

func (m MockNotesStore) NotesRemoveAll(uid string) error {
	return m.Err
}

func (m MockNotesStore) NotesUpdate(uid string, oid string, f io.Reader) error {
	return m.Err
}

func (m MockNotesStore) NotesGetUploadPath(uid string, oid string) string {
	return FakeUploadPath
}

func (m MockNotesStore) NotesCheck(uid string, oid string) (bool, error) {
	return true, m.Err
}
