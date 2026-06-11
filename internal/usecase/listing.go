package usecase

import (
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"better-server/internal/domain"
)

// DirListUseCase handles the logic for listing directory contents.
type DirListUseCase struct {
	repo       domain.FileRepository
	classifier *domain.FileClassifier
}

func NewDirListUseCase(repo domain.FileRepository, classifier *domain.FileClassifier) *DirListUseCase {
	return &DirListUseCase{repo: repo, classifier: classifier}
}

// DirListInput holds parameters for listing a directory.
type DirListInput struct {
	FullPath   string
	RelPath    string
	ShowHidden bool
}

// DirListOutput holds the result of listing a directory.
type DirListOutput struct {
	Files  []domain.FileInfo
	Parent string
}

// Execute lists the contents of a directory and classifies each entry.
func (uc *DirListUseCase) Execute(input DirListInput) (*DirListOutput, error) {
	entries, err := uc.repo.ReadDir(input.FullPath)
	if err != nil {
		return nil, err
	}

	var files []domain.FileInfo
	for _, e := range entries {
		name := e.Name
		isHidden := strings.HasPrefix(name, ".")
		if isHidden && !input.ShowHidden {
			continue
		}

		ext := strings.ToLower(filepath.Ext(name))
		fileType := uc.classifier.Classify(e.IsDir, ext)

		path := "/" + filepath.ToSlash(filepath.Join(input.RelPath, name))
		if input.RelPath == "" {
			path = "/" + name
		}

		files = append(files, domain.FileInfo{
			Name:    name,
			Path:    path,
			Size:    e.Size,
			ModTime: e.ModTime,
			IsDir:   e.IsDir,
			Type:    fileType,
			Ext:     ext,
			Hidden:  isHidden,
		})
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	parent := ""
	if input.RelPath != "" {
		parentDir := filepath.Dir(input.RelPath)
		if parentDir == "." {
			parent = "/"
		} else {
			parent = "/" + filepath.ToSlash(parentDir)
		}
	}

	return &DirListOutput{Files: files, Parent: parent}, nil
}

// BuildBackURL constructs the "back" URL for the player, preserving the hidden param.
func BuildBackURL(relPath, hiddenQuery string) string {
	backUrl := "/"
	if parent := filepath.Dir(relPath); parent != "." {
		backUrl = "/" + filepath.ToSlash(parent)
	}
	if hiddenQuery != "" {
		if strings.Contains(backUrl, "?") {
			backUrl += "&"
		} else {
			backUrl += "?"
		}
		backUrl += "hidden=" + url.QueryEscape(hiddenQuery)
	}
	return backUrl
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
