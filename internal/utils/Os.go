package utils

import "io/fs"

type OS interface {
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
	Open(name string) (fs.File, error)
}
