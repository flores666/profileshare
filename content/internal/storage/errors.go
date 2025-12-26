package storage

import "errors"

var (
	NotFound  = errors.New("url not Found")
	UrlExists = errors.New("url exists")
)
