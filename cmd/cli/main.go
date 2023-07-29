package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ck3g/laterlater/internal/video"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	videoAPI := video.NewYouTubeAPI(apiKey)

	repo, err := video.NewFileRepository("tmp/videos.txt")
	if err != nil {
		fmt.Println("error initializing video repository")
	}

	ids, err := repo.GetAll(context.Background())
	if err != nil {
		fmt.Printf("error fetching video IDs: %w\n", err)
	}

	videos, err := videoAPI.GetInfo(context.Background(), ids)
	if err != nil {
		fmt.Printf("error getting video info: %w\n", err)
	}

	fmt.Println("# Videos")

	for _, v := range videos {
		fmt.Printf(
			"TITLE: %s\nTHUMB: %s\nCHANNEL: %s\nTAGS: %v\nDURATION: %s\n\n",
			v.Title,
			v.ThumbnailURL,
			v.ChannelTitle,
			v.Tags,
			v.Duration,
		)
	}
}
