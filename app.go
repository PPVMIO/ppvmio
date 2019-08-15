package main

import (
	"log"
	"math/rand"
	"net/http"
	"ppvmio/svc"
	"ppvmio/utils"
	"time"
)

var s3BaseUrl = "http://photos-ppvmio-public.s3.amazonaws.com/"

type Layout struct {
	DarkTheme     bool
	Background    bool
	BackgroundUrl string
	Photos        []Photo
}

type Photo struct {
	Url string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	utils.LoadTemplates()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", 301)
	})

	http.HandleFunc("/home", home)
	http.HandleFunc("/photos", photos)
	http.HandleFunc("/moodboard", moodboard)
	http.HandleFunc("/projects", projects)
	http.HandleFunc("/about", about)

	log.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}

func projects(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "projects.html", nil)
}

func moodboard(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "moodboard.html", nil)
}

func about(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "about.html", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	backgrounds := svc.RetrievePhotos("background/")
	l := Layout{true, true, s3BaseUrl + *backgrounds[rand.Intn(len(backgrounds))].Key, nil}
	utils.RenderTemplate(w, "home.html", l)
}

func photos(w http.ResponseWriter, r *http.Request) {
	l := Layout{false, false, "", nil}
	s3Photos := svc.RetrievePhotos("website/")
	var photos []Photo
	for _, item := range s3Photos {
		p := Photo{s3BaseUrl + *item.Key}
		photos = append(photos, p)
	}
	rand.Shuffle(len(photos), func(i, j int) { photos[i], photos[j] = photos[j], photos[i] })
	l.Photos = photos
	utils.RenderTemplate(w, "photos.html", l)
}
