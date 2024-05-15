package project

import (
	"errors"
	"io"
)

type Notes struct {
  Title string
  Data io.ReadCloser
}

func (m *Notes) Validate() error {
  if m.Title == "" {
    return errors.New("Missing Title")
  }

  // TODO: Validate actual contents of data to make sure it is a valid Slate editor value

  return nil
}

// Constructor
func NewNotes(title string, data io.ReadCloser) (Notes, error) {
  n := Notes{
    Title: title + ".stmp",
    Data: data,
  }

  if err := n.Validate(); err != nil {
    return n, nil
  }

  return n, nil

}
