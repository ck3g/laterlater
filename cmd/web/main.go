package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ck3g/laterlater/internal/storage"
	inmemorystorage "github.com/ck3g/laterlater/internal/storage/inmemory"
	"github.com/ck3g/laterlater/internal/video"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

type Application struct {
	Storage  storage.Storage
	YTClient video.YouTubeAPI
}

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

	videoStorage := inmemorystorage.NewInmemoryVideoStorage()
	ytClient := video.NewYouTubeAPI(os.Getenv("API_KEY"))
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
	videos, err := ytClient.GetInfo(context.TODO(), videoIDs)
	if err != nil {
		log.Panic("Error parsing initial video list", err)
	}

	videoStorage.Create(context.TODO(), videos)

	a := Application{
		Storage: storage.Storage{
			Videos: videoStorage,
		},
		YTClient: ytClient,
	}

	viewsEngine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: viewsEngine,
	})
	app.Static("/static", "./static")

	app.Get("/", a.HomeHandler)
	app.Post("/videos", a.CreateVideosHandler)
	app.Delete("/videos/:id", a.DeleteVideoHandler)

	err = app.Listen(":4000")
	if err != nil {
		log.Panic("Error starting a web server: ", err)
	}
}

func (a *Application) HomeHandler(c *fiber.Ctx) error {
	videos, err := a.Storage.Videos.GetAll(c.Context())
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.Render("index", fiber.Map{
		"Videos": videos,
	})
}

// POST handler to create videos
func (a *Application) CreateVideosHandler(c *fiber.Ctx) error {
	videosInput := c.FormValue("videos")

	videoIDs := strings.Split(videosInput, "\n")

	for i, videoID := range videoIDs {
		videoIDs[i] = strings.TrimSpace(videoID)
	}

	videoIDs = video.ParseIDs(videoIDs)
	videos, err := a.YTClient.GetInfo(c.Context(), videoIDs)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	a.Storage.Videos.Create(c.Context(), videos)

	return c.Redirect("/", http.StatusSeeOther)
}

func (a *Application) DeleteVideoHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	a.Storage.Videos.Delete(c.Context(), id)

	return c.JSON(fiber.Map{"result": "ok"})
}
