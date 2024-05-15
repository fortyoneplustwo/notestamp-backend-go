package project


type ProjectStore interface {
  Add(uid int, p Project) error
  Get(uid int, title string) (Project, error)
  List(uid int) ([]string, error)
  Remove(uid int, title string) (Project, error)
}
