package utils

import (
	"io"
	"io/fs"
	"os"
)

type OS interface {
	ReadFile(string) ([]byte, error)
	WriteFile(string, []byte, fs.FileMode) error
	Open(string) (fs.File, error)
	Create(string) (*os.File, error)
	Copy(io.Writer, io.Reader) (int64, error)
	Stat(string) (os.FileInfo, error)
	IsNotExist(error) bool
}
