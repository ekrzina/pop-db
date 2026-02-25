package utils

import (
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
func (RealOS) ReadFile(name string) ([]byte, error) {
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
func (RealOS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

// Open is a os.Open method wrapper
// Parameters:
//   - name: File path to open file from.
//
// Returns:
//   - fs.File: Opened file.
//   - error: Error returned on opening failure.
func (RealOS) Open(name string) (fs.File, error) {
	return os.Open(name)
}
