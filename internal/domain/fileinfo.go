package domain

import (
	"time"
)

// FileInfo represents metadata about a file or directory entry.
type FileInfo struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
	IsDir   bool      `json:"isDir"`
	Type    string    `json:"type"`
	Ext     string    `json:"ext"`
	Hidden  bool      `json:"hidden"`
}

// DownloadRequest represents a request to download multiple paths as a ZIP.
type DownloadRequest struct {
	Paths []string `json:"paths"`
}

// PlayerData holds data for rendering the media player template.
type PlayerData struct {
	Src     string
	Name    string
	Type    string
	BackUrl string
}

// DirListing holds data for rendering the directory listing template.
type DirListing struct {
	Title       string
	CurrentPath string
	Files       []FileInfo
	Parent      string
}
