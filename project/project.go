package project

import (
	"encoding/json"
	"errors"
)

type Project struct {
  Format string `json:"type"`
  Label string `json:"label"`
  Title string `json:"title"`
  Src string `json:"src"`
  Mimetype string `json:"mimetype"`
}


func (p *Project) Validate() error {
  if p.Title == "" {
    return errors.New("Title is empty")
  }

  if p.Format == "" {
    return errors.New("Format is empty")
  }

  if p.Label == "" {
    return errors.New("Label is empty")
  }

  if p.Src != "" && p.Mimetype != "" {
    return errors.New("Src and Mimetype must be mutually exclusive")
  }

  return nil
}


// Constructor
func NewProject(metadata []byte) (Project, error) {
  p := Project{}
  err := json.Unmarshal(metadata, &p)
  if err != nil {
    return p, err
  }

  err = p.Validate()
  if err != nil {
    return p, err
  }

  return p, nil
}

