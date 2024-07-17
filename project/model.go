package project

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
)

type Metadata struct {
	Format   string `json:"type"`
	Label    string `json:"label"`
	Title    string `json:"title"`
	Src      string `json:"src"`
	Mimetype string `json:"mimetype"`
  MediaHash string `json:"mediaHash"`
}

func NewMetadata(j []byte) (Metadata, error) {
	p := Metadata{}
	err := json.Unmarshal(j, &p)
	if err != nil {
		return p, err
	}

	err = p.Validate()
	if err != nil {
		return p, err
	}

	return p, nil
}

func (p *Metadata) Validate() error {
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

func (p *Metadata) MakeMediaKey() string {
  return p.MediaHash + "." + strings.Split(p.Mimetype, "/")[1]
}

func (p *Metadata) MakeNotesKey() string {
  return p.Title + ".stmp"
}

func (p *Metadata) SetMediaHash(f multipart.File) error {
    hasher := sha256.New()
    if _, err := io.Copy(hasher, f); err != nil {
      return err
    }
  _, err := f.Seek(0, io.SeekStart)
    if err != nil {
        return err
    }
    p.MediaHash = fmt.Sprintf("%x", hasher.Sum(nil))

  return nil
}

func (p *Metadata) ToMap() map[string]string {
	m := make(map[string]string)
	m["title"] = p.Title
	m["type"] = p.Format
	m["label"] = p.Label
	m["src"] = p.Src
	m["mimetype"] = p.Mimetype
  return m
}


type Notes struct {
	Title string
	Data  io.ReadCloser
}

func (m *Notes) Validate() error {
	if m.Title == "" {
		return errors.New("Missing Title")
	}

	// TODO: Validate actual contents of data to make sure it is a valid Slate editor value

	return nil
}

func NewNotes(title string, data io.ReadCloser) (Notes, error) {
	n := Notes{
		Title: title + ".stmp",
		Data:  data,
	}

	if err := n.Validate(); err != nil {
		return n, nil
	}

	return n, nil

}


type Media struct {
	Title string
	Data  io.ReadCloser
}

func NewMedia(title string, mimetype string, data io.ReadCloser) (Media, error) {
	m := Media{
		Title: title + "." + strings.Split(mimetype, "/")[1],
		Data:  data,
	}

	err := m.Validate()
	if err != nil {
		return m, err
	}

	return m, nil
}

func (m *Media) Validate() error {
	if m.Title == "" {
		return errors.New("No title provided")
	}

	return nil
}

