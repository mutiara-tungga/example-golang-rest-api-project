package service

import "errors"

var (
	ErrNotFound          = errors.New("user not found")
	ErrDuplicateUsername = errors.New("username already exists")
	ErrWrongPassword     = errors.New("wrong password")
)
