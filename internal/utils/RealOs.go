package utils

import (
	"io"
	"io/fs"
	"os"
)

type RealOS struct{}

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
