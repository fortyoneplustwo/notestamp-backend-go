package project

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type Metadata struct {
	Format   string `json:"type"`
	Label    string `json:"label"`
	Title    string `json:"title"`
	Src      string `json:"src"`
	Mimetype string `json:"mimetype"`
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
	if p.Src != "" && p.Mimetype != "" { // Projects cannot be saved without media
		return errors.New("Src and Mimetype must be mutually exclusive")
	}
	return nil
}

// Constructor
func NewMetadata(metadata []byte) (Metadata, error) {
	p := Metadata{}
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

// Constructor
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

func (m *Media) Validate() error {
	if m.Title == "" {
		return errors.New("No title provided")
	}

	return nil
}

// Constructor
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
