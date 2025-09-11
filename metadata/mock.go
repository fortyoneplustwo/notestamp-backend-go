package metadata

type MockMetadataStore struct {
	Err error
}

var (
	FakeMetadata = Metadata{
		Title:         "testTitle",
		Format:        Audio,
		Src:           "testSrc",
		NotesMimeType: "application/json",
		NotesSize:     500,
	}
	FakeUid = "fakeUid"
)

func (m MockMetadataStore) MetadataList(uid string) ([]Metadata, error) {
	return []Metadata{FakeMetadata}, nil
}

func (m MockMetadataStore) MetadataRemove(uid string, id string) error {
	return m.Err
}

func (m MockMetadataStore) MetadataRemoveAll(uid string) error {
	return m.Err
}

func (m MockMetadataStore) MetadataGet(uid string, oid string) (Metadata, error) {
	return FakeMetadata, m.Err
}

func (m MockMetadataStore) MetadataAdd(uid string, meta Metadata) error {
	return m.Err
}
