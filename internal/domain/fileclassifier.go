package domain

import (
	"fmt"
)

// FileClassifier provides domain logic for classifying files by type.
type FileClassifier struct{}

func NewFileClassifier() *FileClassifier {
	return &FileClassifier{}
}

func (c *FileClassifier) IsImage(ext string) bool {
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp" || ext == ".bmp" || ext == ".svg"
}

func (c *FileClassifier) IsVideo(ext string) bool {
	return ext == ".mp4" || ext == ".webm" || ext == ".ogg" || ext == ".mov" || ext == ".avi"
}

func (c *FileClassifier) IsAudio(ext string) bool {
	return ext == ".mp3" || ext == ".wav" || ext == ".flac" || ext == ".aac" || ext == ".m4a" || ext == ".wma" || ext == ".opus"
}

// Classify returns the file type string: "dir", "image", "video", "audio", or "file".
func (c *FileClassifier) Classify(isDir bool, ext string) string {
	if isDir {
		return "dir"
	}
	if c.IsImage(ext) {
		return "image"
	}
	if c.IsVideo(ext) {
		return "video"
	}
	if c.IsAudio(ext) {
		return "audio"
	}
	return "file"
}

// HumanSize formats bytes into a human-readable string.
func HumanSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
