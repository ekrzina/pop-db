package utils

import (
	"io"
	"io/fs"
	"os"
)

type OS interface {
	Create(string) (*os.File, error)
	Copy(io.Writer, io.Reader) (int64, error)
	IsNotExist(error) bool
	MkdirAll(string, fs.FileMode) error
	Open(string) (fs.File, error)
	ReadDir(name string) ([]os.DirEntry, error)
	ReadFile(string) ([]byte, error)
	Remove(string) error
	Stat(string) (os.FileInfo, error)
	WriteFile(string, []byte, fs.FileMode) error
}
