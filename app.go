package main

import (
	"log"
	"math/rand"
	"net/http"
	"ppvmio/svc"
	"ppvmio/utils"
	"time"

	mobiledetect "github.com/Shaked/gomobiledetect"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3BaseUrl = "http://photos-ppvmio-public.s3.amazonaws.com/"
var backgrounds = svc.RetrieveObjects("background/")

var supremeLogo = Logo{
	Brand:           "Click Me",
	BackgroundColor: "red",
	Color:           "white",
	FontFamily:      "Helvetica,Arial,sans-serif",
	FontWeight:      "bold",
	FontStyle:       "italic",
}

type Layout struct {
	DarkTheme     bool
	Background    bool
	BackgroundUrl string
	Photos        []Photo
	IsMobile      bool
	Logo          Logo
}

type Logo struct {
	Brand           string
	BackgroundColor string
	Color           string
	FontFamily      string
	FontWeight      string
	FontStyle       string
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
	http.HandleFunc("/mood", mood)
	http.HandleFunc("/projects", projects)
	http.HandleFunc("/about", about)

	log.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}

func projects(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "projects.html", nil)
}

func mood(w http.ResponseWriter, r *http.Request) {

	detect := mobiledetect.NewMobileDetect(r, nil)

	var s3Photos = svc.RetrieveObjects("mood/photos/")
	var s3Gifs = svc.RetrieveObjects("mood/gifs/")

	l := Layout{true, true, s3BaseUrl + *s3Gifs[rand.Intn(len(s3Gifs))].Key, nil, detect.IsMobile(), supremeLogo}

	if detect.IsMobile() {
		l.Photos = shufflePhotos(s3Photos)[0:17]
	} else {
		l.Photos = shufflePhotos(s3Photos)[0:29]
	}

	utils.RenderTemplate(w, "mood.html", l)
}

func about(w http.ResponseWriter, r *http.Request) {
	detect := mobiledetect.NewMobileDetect(r, nil)
	l := Layout{true, true, s3BaseUrl + *backgrounds[rand.Intn(len(backgrounds))].Key, nil, detect.IsMobile(), Logo{}}
	utils.RenderTemplate(w, "about.html", l)
}

func home(w http.ResponseWriter, r *http.Request) {
	detect := mobiledetect.NewMobileDetect(r, nil)
	l := Layout{true, true, s3BaseUrl + *backgrounds[rand.Intn(len(backgrounds))].Key, nil, detect.IsMobile(), Logo{}}
	utils.RenderTemplate(w, "home.html", l)
}

func photos(w http.ResponseWriter, r *http.Request) {
	detect := mobiledetect.NewMobileDetect(r, nil)
	l := Layout{false, false, "", nil, detect.IsMobile(), Logo{}}
	s3Photos := svc.RetrieveObjects("website/")
	l.Photos = shufflePhotos(s3Photos)
	utils.RenderTemplate(w, "photos.html", l)
}

func shufflePhotos(s3Photos []*s3.Object) []Photo {
	var photos []Photo
	for _, item := range s3Photos {
		p := Photo{s3BaseUrl + *item.Key}
		photos = append(photos, p)
	}
	rand.Shuffle(len(photos), func(i, j int) { photos[i], photos[j] = photos[j], photos[i] })
	return photos
}
