package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ck3g/laterlater/internal/video"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

var ytClient video.YouTubeAPI

var allVideos = []string{
	"https://www.youtube.com/watch?v=Cs2j-Rjqg94",
	"https://www.youtube.com/watch?v=dJIUxvfSg6A",
	"https://www.youtube.com/watch?v=5EYl1TkJSZY",
	"https://www.youtube.com/watch?v=Lwr3-doAgaI",
	"https://www.youtube.com/watch?v=kWfP4H1qzCk",
	"https://www.youtube.com/watch?v=6FY9urgIjqo",
	"https://www.youtube.com/watch?v=IWDlVSSdKC8",
	"https://www.youtube.com/watch?v=Ztk9d78HgC0",
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ytClient = video.NewYouTubeAPI(os.Getenv("API_KEY"))

	viewsEngine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: viewsEngine,
	})
	app.Static("/static", "./static")

	app.Get("/", homeHandler)
	app.Post("/videos", createVideosHandler)
	app.Delete("/videos/:id", deleteVideoHandler)

	err = app.Listen(":4000")
	if err != nil {
		log.Panic("Error starting a web server: ", err)
	}
}

func homeHandler(c *fiber.Ctx) error {
	videoIDs := video.ParseIDs(allVideos)
	videos, err := ytClient.GetInfo(c.Context(), videoIDs)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.Render("index", fiber.Map{
		"Videos": videos,
	})
}

// POST handler to create videos
func createVideosHandler(c *fiber.Ctx) error {
	videosInput := c.FormValue("videos")

	videos := strings.Split(videosInput, "\n")

	for i, video := range videos {
		videos[i] = strings.TrimSpace(video)
	}

	allVideos = append(allVideos, videos...)

	return c.Redirect("/", http.StatusSeeOther)
}

func deleteVideoHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	for i, video := range allVideos {
		if video == fmt.Sprintf("https://www.youtube.com/watch?v=%s", id) {
			allVideos = append(allVideos[:i], allVideos[i+1:]...)
			break
		}
	}

	return c.JSON(fiber.Map{"result": "ok"})
}
