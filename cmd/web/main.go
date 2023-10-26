package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ck3g/laterlater/internal/video"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

type Application struct {
	Repository video.Repository
	YTClient   video.YouTubeAPI
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	repo, _ := video.NewInMemoryRepository()
	ytClient := video.NewYouTubeAPI(os.Getenv("API_KEY"))

	a := Application{
		Repository: repo,
		YTClient:   ytClient,
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
	allVideos, _ := a.Repository.GetAll(c.Context())
	videoIDs := video.ParseIDs(allVideos)
	videos, err := a.YTClient.GetInfo(c.Context(), videoIDs)
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

	videos := strings.Split(videosInput, "\n")

	for i, video := range videos {
		videos[i] = strings.TrimSpace(video)
	}

	a.Repository.Create(c.Context(), videos)

	return c.Redirect("/", http.StatusSeeOther)
}

func (a *Application) DeleteVideoHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	a.Repository.Delete(c.Context(), id)

	return c.JSON(fiber.Map{"result": "ok"})
}
