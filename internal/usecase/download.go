package usecase

import (
	"archive/zip"
	"bytes"
	"io"
	"path/filepath"

	"better-server/internal/domain"
)

// DownloadUseCase handles the logic for packing files into a ZIP archive.
type DownloadUseCase struct {
	repo domain.FileRepository
}

func NewDownloadUseCase(repo domain.FileRepository) *DownloadUseCase {
	return &DownloadUseCase{repo: repo}
}

// DownloadInput holds the paths to include in the archive.
type DownloadInput struct {
	Paths []string
	Root  string
}

// Execute creates a ZIP archive containing all the requested paths.
func (uc *DownloadUseCase) Execute(input DownloadInput) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for _, p := range input.Paths {
		fullPath := filepath.Join(input.Root, p)
		err := uc.repo.Walk(fullPath, func(path string, entry domain.FileEntry, err error) error {
			if err != nil {
				return err
			}

			rel, _ := filepath.Rel(input.Root, path)
			rel = filepath.ToSlash(rel)

			if entry.IsDir {
				_, err := zipWriter.Create(rel + "/")
				return err
			}

			f, err := zipWriter.Create(rel)
			if err != nil {
				return err
			}

			src, err := uc.repo.Open(path)
			if err != nil {
				return err
			}
			defer src.Close()

			_, err = copyAndClose(f, src)
			return err
		})
		if err != nil {
			zipWriter.Close()
			return nil, err
		}
	}

	zipWriter.Close()
	return buf.Bytes(), nil
}

func copyAndClose(w io.Writer, r io.ReadCloser) (int64, error) {
	n, err := io.Copy(w, r)
	r.Close()
	return n, err
}
