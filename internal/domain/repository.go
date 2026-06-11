package domain

import (
	"io"
	"time"
)

// FileEntry represents a single item read from the file system.
type FileEntry struct {
	Name    string
	Size    int64
	ModTime time.Time
	IsDir   bool
}

// FileRepository defines the interface for accessing the file system.
// This is the port in Clean Architecture / DDD terms.
type FileRepository interface {
	// ReadDir returns all entries in the given directory.
	ReadDir(fullPath string) ([]FileEntry, error)
	// Stat returns file info for a single path.
	Stat(fullPath string) (FileEntry, bool, error)
	// Open opens a file for reading.
	Open(fullPath string) (io.ReadCloser, error)
	// Walk walks a file tree rooted at fullPath, calling walkFn for each file/dir.
	Walk(fullPath string, walkFn func(path string, entry FileEntry, err error) error) error
	// Abs resolves a relative path to an absolute one.
	Abs(path string) (string, error)
}
