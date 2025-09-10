package metadata

import "errors"

var (
	ErrDuplicate = errors.New("project already exists")
	ErrNotExist = errors.New("project does not exits")
)
