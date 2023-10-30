package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ck3g/laterlater/internal/storage"
	mongostorage "github.com/ck3g/laterlater/internal/storage/mongo"
	"github.com/joho/godotenv"
)

const dbName = "laterlater"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := storage.InitMongoDB(os.Getenv("MONGO_URI"))
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	videoStorage, err := mongostorage.NewVideoStorage(client, dbName)
	if err != nil {
		log.Panic("error creating new video storage: ", err)
	}

	videos, err := videoStorage.GetAll(context.Background())
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
