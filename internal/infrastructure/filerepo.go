package infrastructure

import (
	"io"
	"os"
	"path/filepath"

	"better-server/internal/domain"
)

// FileSystemRepo implements domain.FileRepository using the real OS file system.
type FileSystemRepo struct {
}

func NewFileSystemRepo() *FileSystemRepo {
	return &FileSystemRepo{}
}

// ReadDir returns all entries in the given directory.
func (r *FileSystemRepo) ReadDir(fullPath string) ([]domain.FileEntry, error) {
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	result := make([]domain.FileEntry, 0, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		result = append(result, domain.FileEntry{
			Name:    e.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   e.IsDir(),
		})
	}
	return result, nil
}

// Stat returns file info for a single path.
func (r *FileSystemRepo) Stat(fullPath string) (domain.FileEntry, bool, error) {
	fi, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.FileEntry{}, false, nil
		}
		return domain.FileEntry{}, false, err
	}
	return domain.FileEntry{
		Name:    fi.Name(),
		Size:    fi.Size(),
		ModTime: fi.ModTime(),
		IsDir:   fi.IsDir(),
	}, true, nil
}

// Open opens a file for reading.
func (r *FileSystemRepo) Open(fullPath string) (io.ReadCloser, error) {
	return os.Open(fullPath)
}

// Walk walks a file tree rooted at fullPath.
func (r *FileSystemRepo) Walk(fullPath string, walkFn func(path string, entry domain.FileEntry, err error) error) error {
	return filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		var entry domain.FileEntry
		if info != nil {
			entry = domain.FileEntry{
				Name:    info.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime(),
				IsDir:   info.IsDir(),
			}
		}
		return walkFn(path, entry, err)
	})
}

// Abs resolves a relative path to an absolute one.
func (r *FileSystemRepo) Abs(path string) (string, error) {
	return filepath.Abs(path)
}
