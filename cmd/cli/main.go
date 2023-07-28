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
	videos, err := videoAPI.GetInfo(context.Background(), []string{"FNnckb4rg5o", "CK5rLpZk5A8", "WQKPIOvt2Ac"})
	if err != nil {
		fmt.Printf("error getting video info: %w", err)
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
