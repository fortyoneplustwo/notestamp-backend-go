package user

type MockUserStore struct {
	Err  error
	User User
}

var (
	FakeUid         = "fakeUid"
	FakeCredentials = Credentials{
		Email:    "testEmail@gmail.com",
		Password: "fakePassword",
	}
	FakeUser = User{
		Uid:      "fakeUid",
		Email:    "fakeUser",
		Password: "fakePassword",
	}
)

func (m MockUserStore) UserAdd(u Credentials) (string, error) {
	return FakeUid, m.Err
}

func (m MockUserStore) UserGetById(uid string) (User, error) {
	return m.User, m.Err
}

func (m MockUserStore) UserGetByEmail(email string) (User, error) {
	return m.User, m.Err
}

func (m MockUserStore) UserRemove(uid string) error {
	return m.Err
}

