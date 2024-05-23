package project

import (
	"errors"
	"io"
	"strings"
)

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
