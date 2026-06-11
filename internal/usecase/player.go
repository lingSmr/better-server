package usecase

import (
	"path/filepath"
	"strings"

	"better-server/internal/domain"
)

// PlayerUseCase handles the logic for determining media player data.
type PlayerUseCase struct {
	classifier *domain.FileClassifier
}

func NewPlayerUseCase(classifier *domain.FileClassifier) *PlayerUseCase {
	return &PlayerUseCase{classifier: classifier}
}

// PlayerInput holds parameters for building player data.
type PlayerInput struct {
	RelPath string
	Hidden  string
}

// Execute builds PlayerData for a given file path.
func (uc *PlayerUseCase) Execute(input PlayerInput) domain.PlayerData {
	ext := strings.ToLower(filepath.Ext(input.RelPath))
	mediaType := "video"
	if uc.classifier.IsAudio(ext) {
		mediaType = "audio"
	}

	backUrl := BuildBackURL(input.RelPath, input.Hidden)

	return domain.PlayerData{
		Src:     "/" + filepath.ToSlash(input.RelPath),
		Name:    filepath.Base(input.RelPath),
		Type:    mediaType,
		BackUrl: backUrl,
	}
}
