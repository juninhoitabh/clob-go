package shared

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidSide  = errors.New("invalid side")
	ErrInvalidParam = errors.New("invalid parameter")
)
