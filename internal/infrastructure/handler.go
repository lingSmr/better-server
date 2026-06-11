package infrastructure

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"better-server/internal/domain"
	"better-server/internal/usecase"
)

// Handler is the HTTP handler that routes requests to the appropriate use cases.
type Handler struct {
	root       string
	dirListUC  *usecase.DirListUseCase
	downloadUC *usecase.DownloadUseCase
	playerUC   *usecase.PlayerUseCase
	templates  *Templates
	repo       domain.FileRepository
}

// NewHandler creates a new HTTP handler with all its dependencies.
func NewHandler(
	root string,
	dirListUC *usecase.DirListUseCase,
	downloadUC *usecase.DownloadUseCase,
	playerUC *usecase.PlayerUseCase,
	templates *Templates,
	repo domain.FileRepository,
) *Handler {
	return &Handler{
		root:       root,
		dirListUC:  dirListUC,
		downloadUC: downloadUC,
		playerUC:   playerUC,
		templates:  templates,
		repo:       repo,
	}
}

// ServeHTTP implements the http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := filepath.Clean(r.URL.Path)

	// Download route: POST /download
	if r.Method == http.MethodPost && upath == "/download" {
		h.serveDownload(w, r)
		return
	}

	// Player route: /player/some/file.mp4
	if strings.HasPrefix(upath, "/player/") || upath == "/player" {
		h.servePlayer(w, r)
		return
	}

	if upath == "." || upath == "/" {
		h.serveDir(w, r, h.root, "")
		return
	}

	relPath := strings.TrimLeft(upath, "/")
	fullPath := filepath.Join(h.root, relPath)

	_, exists, err := h.repo.Stat(fullPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.NotFound(w, r)
		return
	}

	entry, _, _ := h.repo.Stat(fullPath)
	if entry.IsDir {
		h.serveDir(w, r, fullPath, relPath)
	} else {
		http.ServeFile(w, r, fullPath)
	}
}

func (h *Handler) servePlayer(w http.ResponseWriter, r *http.Request) {
	relPath := strings.TrimPrefix(r.URL.Path, "/player")
	relPath = strings.TrimLeft(relPath, "/")

	if relPath == "" {
		http.NotFound(w, r)
		return
	}

	fullPath := filepath.Join(h.root, relPath)
	_, exists, err := h.repo.Stat(fullPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.NotFound(w, r)
		return
	}

	hidden := r.URL.Query().Get("hidden")
	data := h.playerUC.Execute(usecase.PlayerInput{
		RelPath: relPath,
		Hidden:  hidden,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.templates.Player.Execute(w, data)
}

func (h *Handler) serveDir(w http.ResponseWriter, r *http.Request, fullPath, relPath string) {
	showHidden := false
	if val := r.URL.Query().Get("hidden"); val != "" {
		showHidden, _ = strconv.ParseBool(val)
	}

	output, err := h.dirListUC.Execute(usecase.DirListInput{
		FullPath:   fullPath,
		RelPath:    relPath,
		ShowHidden: showHidden,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := domain.DirListing{
		Title:       "Index of /" + relPath,
		CurrentPath: relPath,
		Files:       output.Files,
		Parent:      output.Parent,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.templates.Index.Execute(w, data)
}

func (h *Handler) serveDownload(w http.ResponseWriter, r *http.Request) {
	var req domain.DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if len(req.Paths) == 0 {
		http.Error(w, "no paths provided", http.StatusBadRequest)
		return
	}

	data, err := h.downloadUC.Execute(usecase.DownloadInput{
		Paths: req.Paths,
		Root:  h.root,
	})
	if err != nil {
		log.Printf("download error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=archive.zip")
	w.Write(data)
}
