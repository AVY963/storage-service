package models

import "time"

type FileMeta struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FileReader interface {
	Read(p []byte) (n int, err error)
	Close() error
}
