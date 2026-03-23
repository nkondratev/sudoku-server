package main

import "errors"

var (
	ErrNotFound error = errors.New("not found")
	ErrEmpty    error = errors.New("empty")
)
