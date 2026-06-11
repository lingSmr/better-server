package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"better-server/internal/domain"
	"better-server/internal/infrastructure"
	"better-server/internal/usecase"
)

func main() {
	port := flag.Int("p", 8080, "порт")
	root := flag.String("d", ".", "директория")
	flag.Parse()

	// --- Infrastructure Layer ---
	repo := infrastructure.NewFileSystemRepo()
	templates := infrastructure.NewTemplates()

	absRoot, err := repo.Abs(*root)
	if err != nil {
		log.Fatalf("failed to resolve root path: %v", err)
	}

	// --- Domain Layer ---
	classifier := domain.NewFileClassifier()

	// --- Use Case Layer ---
	dirListUC := usecase.NewDirListUseCase(repo, classifier)
	downloadUC := usecase.NewDownloadUseCase(repo)
	playerUC := usecase.NewPlayerUseCase(classifier)

	// --- HTTP Handler (Infrastructure) ---
	handler := infrastructure.NewHandler(absRoot, dirListUC, downloadUC, playerUC, templates, repo)

	log.Printf("Сервер запущен → http://localhost:%d", *port)
	log.Printf("Директория: %s", absRoot)

	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
