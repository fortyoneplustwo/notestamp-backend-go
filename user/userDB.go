package user

import (
  "database/sql"
  _ "github.com/lib/pq"
  "errors"
)

type UserDB struct {
  db *sql.DB
}


// errors
var ErrUserNotFound = errors.New("user not found")


// Contructor
func NewUserDB(db *sql.DB) *UserDB {
  return &UserDB{db: db}
}


// Implements UserStore interface:
func (u *UserDB) Add(user User) error {
  stmt, err := u.db.Prepare("INSERT INTO users (email, password) VALUES ($1, $2)")
  if err != nil {
    return err
  }
  defer stmt.Close()

  if _, err := stmt.Exec(user.Email, user.Password); err != nil {
    return err
  }

  return nil
}


func (u *UserDB) Exists(email string) (bool, error) {
  stmt := "SELECT * FROM users WHERE email = $1"
  rows, err := u.db.Query(stmt, email)
  if err != nil {
    return false, err
  }
  defer rows.Close()

  if rows.Next() {
    return true, nil
  } else {
    return false, nil
  }
}


func (u *UserDB) Get(email string) (User, error) {
  stmt := "SELECT id, email, password, directory FROM users WHERE email = $1"
  row := u.db.QueryRow(stmt, email)

  user := User{}
  if err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Directory); err != nil {
    return user, err
  }

  return user, nil
}


func (u *UserDB) Remove(id int) error {
  stmt, err := u.db.Prepare("DELETE FROM users WHERE id = $1")
  if err != nil {
    return err
  }
  defer stmt.Close()

  if _, err := stmt.Exec(id); err != nil {
    if err == sql.ErrNoRows {
      return ErrUserNotFound
    }
    return err
  }

  return nil
}


func (u *UserDB) UpdateDir(id int, newDir []string) error {
  stmt, err := u.db.Prepare("UPDATE users SET directory = $1 WHERE id = $2")
  if err != nil {
    return err
  }
  defer stmt.Close()

  if _, err = stmt.Exec(newDir); err != nil {
    return err
  }

  return nil
}

