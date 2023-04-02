package model

import (
	"errors"
)

var (
	NotFoundError      = errors.New("no result in database")
	DuplicateNickError = errors.New("the provided nickname is already taken")
	DuplicateMailError = errors.New("the provided mail address is already taken")
)
