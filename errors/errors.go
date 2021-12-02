package errors

import "errors"

var (
	ErrInvalidFSON              = errors.New("invalid FSON")
	ErrNotAFolder               = errors.New("not a folder")
	ErrNotAnObject              = errors.New("not an object")
	ErrNotAnArray               = errors.New("not an array")
	ErrIndexOutOfBounds         = errors.New("index out of bounds")
	ErrCannotAccessFileChildren = errors.New("cannot access file children")
	ErrNoParent                 = errors.New("no parent")
)
