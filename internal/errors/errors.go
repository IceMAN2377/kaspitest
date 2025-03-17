package errs

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("user(s) not found")
	ErrAlreadyExists = errors.New("user with such data already exists")
	ErrIncorrectData = errors.New("incorrect data")
)
