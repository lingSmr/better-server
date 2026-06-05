package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	port = flag.Int("p", 8080, "порт")
	root = flag.String("d", ".", "директория")
)

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

var tmpl = template.Must(template.New("index").Funcs(template.FuncMap{
	"humanSize": humanSize,
	"formatTime": func(t time.Time) string {
		return t.Format("02.01.2006 15:04")
	},
}).Parse(`<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css">
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; }
        .card { transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1); }
        .card:hover { transform: translateY(-4px); box-shadow: 0 20px 25px -5px rgb(0 0 0 / 0.1); }
        .preview { height: 180px; }
    </style>
</head>
<body class="bg-zinc-950 text-zinc-200">
    <div class="max-w-7xl mx-auto p-6">
        <div class="flex items-center justify-between mb-8">
            <div>
                <h1 class="text-3xl font-bold flex items-center gap-3">
                    <i class="fas fa-folder text-amber-400"></i>
                    {{.Title}}
                </h1>
                <p class="text-zinc-500 mt-1">{{.CurrentPath}}</p>
            </div>

            <div class="flex items-center gap-4">
                <button onclick="toggleHidden()" id="hiddenBtn"
                        class="flex items-center gap-2 px-4 py-2 rounded-xl bg-zinc-800 hover:bg-zinc-700 transition-colors">
                    <i class="fas fa-eye"></i>
                    <span id="hiddenText">Скрытые: выкл</span>
                </button>
            </div>
        </div>

        {{if .Parent}}
        <a href="{{.Parent}}" class="inline-flex items-center gap-2 text-blue-400 hover:text-blue-300 mb-6 text-lg">
            <i class="fas fa-arrow-left"></i> На уровень выше
        </a>
        {{end}}

        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
            {{range .Files}}
            <a href="{{.Path}}" class="card group bg-zinc-900 rounded-2xl overflow-hidden border border-zinc-800 hover:border-zinc-700 flex flex-col h-full">
                <div class="preview flex items-center justify-center bg-zinc-950 relative overflow-hidden">
                    {{if eq .Type "dir"}}
                        <div class="text-7xl text-amber-400/80 group-hover:scale-110 transition-transform">📁</div>
                    {{else if eq .Type "image"}}
                        <img src="{{.Path}}" class="w-full h-full object-cover" alt="{{.Name}}" loading="lazy">
                    {{else if eq .Type "video"}}
                        <video src="{{.Path}}" class="w-full h-full object-cover" muted preload="metadata"></video>
                        <div class="absolute inset-0 flex items-center justify-center">
                            <i class="fas fa-play text-white/70 text-4xl"></i>
                        </div>
                    {{else}}
                        <div class="text-6xl text-zinc-600 group-hover:text-zinc-500 transition-colors">
                            {{if eq .Ext ".pdf"}}📕{{else if or (eq .Ext ".go") (eq .Ext ".js") (eq .Ext ".py")}}💾{{else}}📄{{end}}
                        </div>
                    {{end}}
                </div>
                
                <div class="p-4 flex-1 flex flex-col">
                    <div class="font-medium text-zinc-100 line-clamp-2 group-hover:text-white transition-colors flex items-center gap-2">
                        {{if .Hidden}}<span class="text-amber-400 text-xs">•</span>{{end}}
                        {{.Name}}
                    </div>
                    <div class="mt-auto pt-3 text-xs text-zinc-500 flex justify-between">
                        <span>{{if not .IsDir}}{{.Size | humanSize}}{{else}}—{{end}}</span>
                        <span>{{.ModTime | formatTime}}</span>
                    </div>
                </div>
            </a>
            {{end}}
        </div>
    </div>

    <script>
        let showHidden = localStorage.getItem('showHidden') === 'true';

        function updateButton() {
            const btn = document.getElementById('hiddenBtn');
            const text = document.getElementById('hiddenText');
            if (showHidden) {
                btn.classList.add('bg-emerald-600', 'hover:bg-emerald-500');
                text.textContent = 'Скрытые: вкл';
            } else {
                btn.classList.remove('bg-emerald-600', 'hover:bg-emerald-500');
                text.textContent = 'Скрытые: выкл';
            }
        }

        function toggleHidden() {
            showHidden = !showHidden;
            localStorage.setItem('showHidden', showHidden);
            updateButton();
            window.location.search = showHidden ? '?hidden=true' : '';
        }

        // Применяем сохранённое состояние и query-параметр
        window.onload = function() {
            const urlParams = new URLSearchParams(window.location.search);
            if (urlParams.has('hidden')) {
                showHidden = urlParams.get('hidden') === 'true';
                localStorage.setItem('showHidden', showHidden);
            }
            updateButton();
        };
    </script>
</body>
</html>`))


func humanSize(bytes int64) string {
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

type handler struct {
	root string
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := filepath.Clean(r.URL.Path)
	if upath == "." || upath == "/" {
		h.serveDir(w, r, h.root, "")
		return
	}

	relPath := strings.TrimLeft(upath, "/")
	fullPath := filepath.Join(h.root, relPath)

	fi, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if fi.IsDir() {
		h.serveDir(w, r, fullPath, relPath)
	} else {
		http.ServeFile(w, r, fullPath)
	}
}

func (h *handler) serveDir(w http.ResponseWriter, r *http.Request, fullPath, relPath string) {
	showHidden := false
	if val := r.URL.Query().Get("hidden"); val != "" {
		showHidden, _ = strconv.ParseBool(val)
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var files []FileInfo
	for _, e := range entries {
		name := e.Name()
		isHidden := strings.HasPrefix(name, ".")

		if isHidden && !showHidden {
			continue
		}

		info, _ := e.Info()
		ext := strings.ToLower(filepath.Ext(name))
		fileType := "file"

		if e.IsDir() {
			fileType = "dir"
		} else if isImage(ext) {
			fileType = "image"
		} else if isVideo(ext) {
			fileType = "video"
		}

		path := "/" + filepath.ToSlash(filepath.Join(relPath, name))
		if relPath == "" {
			path = "/" + name
		}

		files = append(files, FileInfo{
			Name:    name,
			Path:    path,
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   e.IsDir(),
			Type:    fileType,
			Ext:     ext,
			Hidden:  isHidden,
		})
	}

	// Сортировка
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	parent := ""
	if relPath != "" {
		parentDir := filepath.Dir(relPath)
		if parentDir == "." {
			parent = "/"
		} else {
			parent = "/" + filepath.ToSlash(parentDir)
		}
	}

	data := struct {
		Title       string
		CurrentPath string
		Files       []FileInfo
		Parent      string
	}{
		Title:       "Index of /" + relPath,
		CurrentPath: relPath,
		Files:       files,
		Parent:      parent,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

func isImage(ext string) bool {
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp" || ext == ".bmp" || ext == ".svg"
}

func isVideo(ext string) bool {
	return ext == ".mp4" || ext == ".webm" || ext == ".ogg" || ext == ".mov" || ext == ".avi"
}

func main() {
	flag.Parse()

	absRoot, _ := filepath.Abs(*root)

	log.Printf("Сервер запущен → http://localhost:%d", *port)
	log.Printf("Директория: %s", absRoot)

	http.Handle("/", &handler{root: absRoot})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
