package infrastructure

import (
	"html/template"
	"time"

	"better-server/internal/domain"
)

// Templates holds all parsed HTML templates.
type Templates struct {
	Index  *template.Template
	Player *template.Template
}

// NewTemplates parses and returns all templates.
func NewTemplates() *Templates {
	funcMap := template.FuncMap{
		"humanSize": domain.HumanSize,
		"formatTime": func(t time.Time) string {
			return t.Format("02.01.2006 15:04")
		},
	}

	playerTmpl := template.Must(template.New("player").Parse(playerTemplateHTML))
	indexTmpl := template.Must(template.New("index").Funcs(funcMap).Parse(indexTemplateHTML))

	return &Templates{
		Index:  indexTmpl,
		Player: playerTmpl,
	}
}

const playerTemplateHTML = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Плеер — {{.Name}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css">
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; background: #000; }
        .player-container {
            position: fixed;
            inset: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            background: #000;
        }
        .player-container video, .player-container audio {
            max-width: 100%;
            max-height: 100%;
        }
        .player-container audio { width: 100%; max-width: 600px; }
        .controls {
            position: fixed;
            bottom: 0;
            left: 0;
            right: 0;
            background: linear-gradient(transparent, rgba(0,0,0,0.9));
            padding: 40px 20px 20px;
            transition: opacity 0.3s;
        }
        .controls-hidden { opacity: 0; pointer-events: none; }
        .progress-bar {
            position: relative;
            width: 100%;
            height: 6px;
            background: rgba(255,255,255,0.2);
            border-radius: 3px;
            cursor: pointer;
            margin-bottom: 12px;
            transition: height 0.15s;
        }
        .progress-bar:hover { height: 10px; }
        .progress-bar .progress-fill {
            height: 100%;
            background: #22c55e;
            border-radius: 3px;
            position: relative;
            transition: width 0.1s linear;
        }
        .progress-bar .progress-fill::after {
            content: '';
            position: absolute;
            right: -6px;
            top: 50%;
            transform: translateY(-50%);
            width: 12px;
            height: 12px;
            background: #22c55e;
            border-radius: 50%;
            opacity: 0;
            transition: opacity 0.15s;
        }
        .progress-bar:hover .progress-fill::after { opacity: 1; }
        .controls-row {
            display: flex;
            align-items: center;
            gap: 16px;
            color: white;
        }
        .controls-row button {
            background: none;
            border: none;
            color: white;
            cursor: pointer;
            font-size: 20px;
            padding: 8px;
            border-radius: 50%;
            transition: background 0.2s;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .controls-row button:hover { background: rgba(255,255,255,0.15); }
        .controls-row .time { font-size: 14px; font-variant-numeric: tabular-nums; min-width: 45px; }
        .controls-row .time-separator { color: rgba(255,255,255,0.5); }
        .controls-row .time-duration { color: rgba(255,255,255,0.6); }
        .volume-container {
            display: flex;
            align-items: center;
            gap: 8px;
            margin-left: auto;
        }
        .volume-container input[type="range"] {
            width: 80px;
            height: 4px;
            -webkit-appearance: none;
            appearance: none;
            background: rgba(255,255,255,0.3);
            border-radius: 2px;
            outline: none;
        }
        .volume-container input[type="range"]::-webkit-slider-thumb {
            -webkit-appearance: none;
            width: 12px;
            height: 12px;
            background: white;
            border-radius: 50%;
            cursor: pointer;
        }
        .volume-container input[type="range"]::-moz-range-thumb {
            width: 12px;
            height: 12px;
            background: white;
            border-radius: 50%;
            border: none;
            cursor: pointer;
        }
        .speed-btn {
            font-size: 12px !important;
            font-weight: 600;
            padding: 4px 8px !important;
            border-radius: 6px !important;
            background: rgba(255,255,255,0.1);
        }
        .speed-btn:hover { background: rgba(255,255,255,0.25); }
        .back-btn {
            position: fixed;
            top: 16px;
            left: 16px;
            z-index: 50;
            background: rgba(0,0,0,0.6);
            color: white;
            border: none;
            padding: 10px 16px;
            border-radius: 10px;
            cursor: pointer;
            font-size: 16px;
            display: flex;
            align-items: center;
            gap: 8px;
            transition: background 0.2s;
        }
        .back-btn:hover { background: rgba(0,0,0,0.9); }
        @media (max-width: 640px) {
            .volume-container input[type="range"] { width: 50px; }
            .controls-row { gap: 10px; }
            .speed-btn { display: none; }
        }
    </style>
</head>
<body>
    <a href="{{.BackUrl}}" class="back-btn"><i class="fas fa-arrow-left"></i> Назад</a>
    <div class="player-container" id="playerContainer">
        {{if eq .Type "audio"}}
        <audio id="media" src="{{.Src}}"></audio>
        {{else}}
        <video id="media" src="{{.Src}}"></video>
        {{end}}
    </div>
    <div class="controls" id="controls">
        <div class="progress-bar" id="progressBar" onclick="seek(event)">
            <div class="progress-fill" id="progressFill" style="width: 0%"></div>
        </div>
        <div class="controls-row">
            <button onclick="togglePlay()" id="playBtn" title="Воспроизведение/Пауза">
                <i class="fas fa-play"></i>
            </button>
            <span class="time" id="currentTime">0:00</span>
            <span class="time-separator">/</span>
            <span class="time-duration" id="duration">0:00</span>
            <button onclick="skipBackward()" title="Назад 10с">
                <i class="fas fa-backward"></i>
            </button>
            <button onclick="skipForward()" title="Вперед 10с">
                <i class="fas fa-forward"></i>
            </button>
            <div class="volume-container">
                <button onclick="toggleMute()" id="muteBtn" title="Выключить звук">
                    <i class="fas fa-volume-up"></i>
                </button>
                <input type="range" id="volumeSlider" min="0" max="1" step="0.05" value="1" oninput="setVolume(this.value)">
            </div>
            <button onclick="cycleSpeed()" class="speed-btn" id="speedBtn" title="Скорость">1x</button>
            <button onclick="toggleFullscreen()" id="fullscreenBtn" title="На весь экран">
                <i class="fas fa-expand"></i>
            </button>
        </div>
    </div>
    <script>
        const media = document.getElementById('media');
        const playBtn = document.getElementById('playBtn');
        const progressFill = document.getElementById('progressFill');
        const progressBar = document.getElementById('progressBar');
        const currentTimeEl = document.getElementById('currentTime');
        const durationEl = document.getElementById('duration');
        const muteBtn = document.getElementById('muteBtn');
        const volumeSlider = document.getElementById('volumeSlider');
        const speedBtn = document.getElementById('speedBtn');
        const controls = document.getElementById('controls');
        const playerContainer = document.getElementById('playerContainer');

        let hideTimeout = null;
        let isSeeking = false;

        function formatTime(s) {
            if (isNaN(s) || s === Infinity) return '0:00';
            const m = Math.floor(s / 60);
            const sec = Math.floor(s % 60);
            return m + ':' + (sec < 10 ? '0' : '') + sec;
        }

        function togglePlay() {
            if (media.paused) {
                media.play();
                playBtn.innerHTML = '<i class="fas fa-pause"></i>';
            } else {
                media.pause();
                playBtn.innerHTML = '<i class="fas fa-play"></i>';
            }
        }

        function seek(e) {
            const rect = progressBar.getBoundingClientRect();
            const pos = (e.clientX - rect.left) / rect.width;
            media.currentTime = pos * media.duration;
        }

        function skipBackward() { media.currentTime = Math.max(0, media.currentTime - 10); }
        function skipForward() { media.currentTime = Math.min(media.duration, media.currentTime + 10); }

        function setVolume(val) {
            media.volume = val;
            media.muted = false;
            updateMuteIcon();
        }

        function toggleMute() {
            media.muted = !media.muted;
            updateMuteIcon();
        }

        function updateMuteIcon() {
            if (media.muted || media.volume === 0) {
                muteBtn.innerHTML = '<i class="fas fa-volume-mute"></i>';
                volumeSlider.value = 0;
            } else if (media.volume < 0.5) {
                muteBtn.innerHTML = '<i class="fas fa-volume-down"></i>';
                volumeSlider.value = media.volume;
            } else {
                muteBtn.innerHTML = '<i class="fas fa-volume-up"></i>';
                volumeSlider.value = media.volume;
            }
        }

        const speeds = [0.5, 0.75, 1, 1.25, 1.5, 2];
        let speedIndex = 2;
        function cycleSpeed() {
            speedIndex = (speedIndex + 1) % speeds.length;
            media.playbackRate = speeds[speedIndex];
            speedBtn.textContent = speeds[speedIndex] + 'x';
        }

        function toggleFullscreen() {
            if (!document.fullscreenElement) {
                playerContainer.requestFullscreen();
                fullscreenBtn.innerHTML = '<i class="fas fa-compress"></i>';
            } else {
                document.exitFullscreen();
                fullscreenBtn.innerHTML = '<i class="fas fa-expand"></i>';
            }
        }

        function showControls() {
            controls.classList.remove('controls-hidden');
            clearTimeout(hideTimeout);
            if (!media.paused) {
                hideTimeout = setTimeout(() => controls.classList.add('controls-hidden'), 3000);
            }
        }

        // Events
        media.addEventListener('loadedmetadata', () => {
            durationEl.textContent = formatTime(media.duration);
        });

        media.addEventListener('timeupdate', () => {
            if (!isSeeking && media.duration) {
                progressFill.style.width = (media.currentTime / media.duration * 100) + '%';
                currentTimeEl.textContent = formatTime(media.currentTime);
            }
        });

        media.addEventListener('play', () => {
            playBtn.innerHTML = '<i class="fas fa-pause"></i>';
            showControls();
        });

        media.addEventListener('pause', () => {
            playBtn.innerHTML = '<i class="fas fa-play"></i>';
            controls.classList.remove('controls-hidden');
            clearTimeout(hideTimeout);
        });

        media.addEventListener('volumechange', updateMuteIcon);

        document.addEventListener('keydown', (e) => {
            if (e.key === ' ' || e.key === 'Space' || e.code === 'Space') {
                e.preventDefault();
                togglePlay();
            }
            if (e.key === 'ArrowLeft') skipBackward();
            if (e.key === 'ArrowRight') skipForward();
            if (e.key === 'f' || e.key === 'F') toggleFullscreen();
            if (e.key === 'm' || e.key === 'M') toggleMute();
        });

        playerContainer.addEventListener('mousemove', showControls);
        controls.addEventListener('mousemove', showControls);

        // Progress bar drag support
        progressBar.addEventListener('mousedown', (e) => {
            isSeeking = true;
            seek(e);
        });
        document.addEventListener('mousemove', (e) => {
            if (isSeeking) {
                const rect = progressBar.getBoundingClientRect();
                let pos = (e.clientX - rect.left) / rect.width;
                pos = Math.max(0, Math.min(1, pos));
                progressFill.style.width = (pos * 100) + '%';
                currentTimeEl.textContent = formatTime(pos * media.duration);
            }
        });
        document.addEventListener('mouseup', (e) => {
            if (isSeeking) {
                seek(e);
                isSeeking = false;
            }
        });

        // Autoplay
        media.play().catch(() => {});

        // Audio: hide fullscreen button
        if (media.tagName === 'AUDIO') {
            document.getElementById('fullscreenBtn').style.display = 'none';
            playerContainer.style.alignItems = 'center';
            playerContainer.style.padding = '40px';
        }
    </script>
</body>
</html>`

const indexTemplateHTML = `<!DOCTYPE html>
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
        .file-checkbox { position: relative; z-index: 10; cursor: pointer; }
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

        <div class="flex justify-end mb-4" id="downloadBar" style="display:none;">
            <button onclick="downloadSelected()"
                    class="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-blue-600 hover:bg-blue-500 transition-colors text-white font-medium shadow-lg">
                <i class="fas fa-download"></i>
                Скачать выбранное (<span id="selectedCount">0</span>)
            </button>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
            {{range .Files}}
            <div class="card group bg-zinc-900 rounded-2xl overflow-hidden border border-zinc-800 hover:border-zinc-700 flex flex-col h-full relative">
                <label class="absolute top-2 left-2 z-10 w-5 h-5 cursor-pointer" onclick="event.stopPropagation()">
                    <input type="checkbox" class="file-checkbox w-5 h-5 accent-blue-500" data-path="{{.Path}}" onchange="updateDownloadBar()">
                </label>
                <a href="{{if or (eq .Type "video") (eq .Type "audio")}}/player{{.Path}}{{else}}{{.Path}}{{end}}" class="flex flex-col h-full">
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
                        {{else if eq .Type "audio"}}
                            <div class="text-6xl text-zinc-600">
                                <i class="fas fa-music"></i>
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
            </div>
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

        function updateDownloadBar() {
            const checkboxes = document.querySelectorAll('.file-checkbox:checked');
            const bar = document.getElementById('downloadBar');
            const count = document.getElementById('selectedCount');
            if (checkboxes.length > 0) {
                bar.style.display = 'flex';
                count.textContent = checkboxes.length;
            } else {
                bar.style.display = 'none';
            }
        }

        function downloadSelected() {
            const checkboxes = document.querySelectorAll('.file-checkbox:checked');
            if (checkboxes.length === 0) return;
            const paths = Array.from(checkboxes).map(cb => cb.dataset.path);
            fetch('/download', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({paths: paths})
            }).then(res => {
                if (!res.ok) throw new Error('Ошибка');
                return res.blob();
            }).then(blob => {
                const url = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = 'archive.zip';
                document.body.appendChild(a);
                a.click();
                a.remove();
                URL.revokeObjectURL(url);
            }).catch(err => {
                alert('Ошибка при скачивании: ' + err.message);
            });
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
</html>`
