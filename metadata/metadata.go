package metadata

import (
  "errors"
)

type Metadata struct {
  Title string `json:"title"`
  Format string `json:"type"`
  Label string `json:"label"`
  Src string `json:"src"`
  Mimetype string `json:"mimetype"`
}

func (m *Metadata) Validate() error {
  if m.Title == "" {
    return errors.New("Title is empty")
  }

  if m.Format == "" {
    return errors.New("Format is empty")
  }

  if m.Label == "" {
    return errors.New("Label is empty")
  }

  return nil
}
