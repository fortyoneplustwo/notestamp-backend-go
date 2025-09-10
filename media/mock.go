package media

type MockMediaStore struct {
	Err error
}

var (
	FakeUploadPath = "fakeUploadPath"
)

func (m MockMediaStore) MediaRemove(uid string, oid string) error {
	return m.Err
}

func (m MockMediaStore) MediaRemoveAll(uid string) error {
	return m.Err
}

func (m MockMediaStore) MediaGetUploadPath(uid string, oid string) string {
	return FakeUploadPath
}

func (m MockMediaStore) MediaCheck(uid string, oid string) (bool, error) {
	return true, m.Err
}
