package shared

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidParam  = errors.New("invalid parameter")
	ErrAlreadyExists = errors.New("already exists")
)
