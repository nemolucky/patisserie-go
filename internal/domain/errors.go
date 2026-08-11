package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrExpired = errors.New("time expired")