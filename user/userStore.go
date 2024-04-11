package user

type UserStore interface {
  Add(u User) error
  Exists(email string) (bool, error)
  Get(email string) (User, error)
  Remove(id int) error
}
