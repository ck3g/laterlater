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

	repo, err := video.NewInMemoryRepository()
	if err != nil {
		log.Fatalln("error initializing video repository")
	}

	ulrs, err := repo.GetAll(context.Background())
	if err != nil {
		log.Fatalf("error fetching video urls: %v\n", err)
	}

	ids := video.ParseIDs(ulrs)

	if len(ids) > 0 {
		videos, err := videoAPI.GetInfo(context.Background(), ids)
		if err != nil {
			log.Fatalf("error getting video info: %v\n", err)
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
	} else {
		fmt.Println("No videos found!")
	}
}
