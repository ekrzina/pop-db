package utils

import (
	"io"
	"io/fs"
	"os"
)

// Interface implementation
type RealOS struct{}

// Interface compile-time guard
var _ OS = (*RealOS)(nil)

// ReadFile is a os.ReadFile method wrapper
// Parameters:
//   - name: File path to read.
//
// Returns:
//   - []byte - Read bytes from file path.
//   - error: Error is return on file read fail.
func (r *RealOS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// WriteFile is a os.WriteFile method wrapper
// Parameters:
//   - name: File path to write to.
//   - data: Data to write to specified path.
//   - perm: Permissions to set for the written file.
//
// Returns:
//   - error: Error on file writing fail.
func (r *RealOS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

// MkdirAll is a os.MkdirAll method wrapper
// Parameters:
//   - name: File path to write to.
//   - perm: Permissions to set for the written file.
//
// Returns:
//   - error: Error on making dir fail.
func (r *RealOS) MkdirAll(name string, perm fs.FileMode) error {
	return os.Mkdir(name, perm)
}

// ReadDir is a os.ReadDir method wrapper
// Parameters:
//   - name: File path to read.
//
// Returns:
//   - []os.DirEntry: Directory entry interface for reading direcory.
//   - error: Error returned on directory reading fail.
func (r *RealOS) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

// Remove is a os.Remove method wrapper
// Parameters:
//   - name: File path to remove.
//
// Returns:
//   - error: Returns error on removal fail.
func (r *RealOS) Remove(name string) error {
	return os.Remove(name)
}

// Open is a os.Open method wrapper
// Parameters:
//   - name: File path to open file from.
//
// Returns:
//   - fs.File: Opened file.
//   - error: Error returned on opening failure.
func (r *RealOS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

// Open is a os.Create method wrapper
// Parameters:
//   - name: File path to open file from.
//
// Returns:
//   - os.File: An open file descriptor.
//   - error: Error returned on creation failure.
func (r *RealOS) Create(name string) (*os.File, error) {
	return os.Create(name)
}

// Stat is a os.Stat method wrapper
// Parameters:
//   - name: File path to open file from.
//
// Returns:
//   - os.FileInfo: File descriptor describing the named file.
//   - error: Error returned on stat failure.
func (r *RealOS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// Copy is a os.IsNotExist method wrapper
// Parameters:
//   - err: Error to check if OS directory exists.
//
// Returns:
//   - bool: Indication of weather file exists.
func (r *RealOS) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

// Copy is a io.Copy method wrapper
// Parameters:
//   - dst: Writer interface that writes bytes to an underlying data stream.
//   - src: Reader interface that reads up to len(p) bytes into p.
//
// Returns:
//   - int64: Length of data written.
//   - error: Error returned on copy failure.
func (r *RealOS) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}
