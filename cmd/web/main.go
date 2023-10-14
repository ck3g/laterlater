package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/ck3g/laterlater/internal/video"
	"github.com/joho/godotenv"
)

var ytClient video.YouTubeAPI

var allVideos = []string{
	"https://www.youtube.com/watch?v=Cs2j-Rjqg94",
	"https://www.youtube.com/watch?v=dJIUxvfSg6A",
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ytClient = video.NewYouTubeAPI(os.Getenv("API_KEY"))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/videos", createVideosHandler)

	fmt.Println("Starting a web server on port 4000...")
	err = http.ListenAndServe(":4000", nil)
	if err != nil {
		log.Panic("Error starting a web server: ", err)
	}
}

type HomePageData struct {
	Videos []video.Video
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, "Something went wrong. Please try again later.", http.StatusInternalServerError)
		return
	}

	videoIDs := video.ParseIDs(allVideos)
	videos, err := ytClient.GetInfo(context.Background(), videoIDs)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong. Please try again later.", http.StatusInternalServerError)
		return
	}

	data := HomePageData{
		Videos: videos,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Something went wrong. Please try again later.", http.StatusInternalServerError)
		return
	}
}

// POST handler to create videos
func createVideosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported.", http.StatusBadRequest)
		return
	}

	// Read value from "videos" body param
	err := r.ParseForm()
	videos := r.FormValue("videos")
	if err != nil {
		http.Error(w, "Something went wrong. Please try again later.", http.StatusInternalServerError)
		return
	}

	// add received videos to allVideos slice
	allVideos = append(allVideos, strings.Split(videos, ",")...)

	// redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
