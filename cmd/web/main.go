package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ck3g/laterlater/internal/video"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	client := connectToMongoDB(os.Getenv("MONGO_URI"))
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

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

func connectToMongoDB(uri string) *mongo.Client {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}

	log.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client
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
