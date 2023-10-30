package main

import (
	"context"
	"fmt"
	"log"
	"os"

	inmemorystorage "github.com/ck3g/laterlater/internal/storage/inmemory"
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

	videoIDs := []string{
		"https://www.youtube.com/watch?v=i7ABlHngi1Q",
		"https://www.youtube.com/watch?v=Cs2j-Rjqg94",
		"https://www.youtube.com/watch?v=dJIUxvfSg6A",
		"https://www.youtube.com/watch?v=5EYl1TkJSZY",
		"https://www.youtube.com/watch?v=Lwr3-doAgaI",
		"https://www.youtube.com/watch?v=kWfP4H1qzCk",
		"https://www.youtube.com/watch?v=6FY9urgIjqo",
		"https://www.youtube.com/watch?v=IWDlVSSdKC8",
		"https://www.youtube.com/watch?v=Ztk9d78HgC0",
	}
	videoIDs = video.ParseIDs(videoIDs)
	videos, err := videoAPI.GetInfo(context.TODO(), videoIDs)
	if err != nil {
		log.Panic("Error parsing initial video list", err)
	}

	videoStorage := inmemorystorage.NewInmemoryVideoStorage()
	videoStorage.Create(context.TODO(), videos)

	videos, err = videoStorage.GetAll(context.Background())
	if err != nil {
		log.Fatalf("error fetching video urls: %v\n", err)
	}

	if len(videos) > 0 {
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
