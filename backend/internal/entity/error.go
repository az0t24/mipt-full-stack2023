package entity

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrVisitNotFound = errors.New("visit not found")
var ErrInvalidEmail = errors.New("invalid email")
var ErrInvalidPassword = errors.New("invalid password")
var ErrEmptyAuthHeader = errors.New("empty auth header")
var ErrInvalidAuthHeader = errors.New("invalid auth header")
var ErrUserExists = errors.New("user with this email exists")
