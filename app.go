package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"ppvmio/svc"
	"ppvmio/utils"
	"strings"
	"time"

	mobiledetect "github.com/Shaked/gomobiledetect"
	log "github.com/sirupsen/logrus"
)

var cloudFrontBaseURL = "https://du6xmiczrsmmq.cloudfront.net/"
var moodboardPagePhotos []Photo
var photosPagePhotos []Photo
var moodboardBackgroundGifs []string
var backgroundImages []string
var aboutBackgroundGifs []string

func init() {
	rand.Seed(time.Now().UnixNano())
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
}

func main() {

	go refreshBackgroundImages()
	go refreshPhotosPagePhotos()
	go refreshMoodboardPagePhotos()
	go refreshMoodboardBackgroundGifs()
	go refreshAboutGifs()

	utils.LoadTemplates()

	http.HandleFunc("/", home)
	http.HandleFunc("/home", home)
	http.HandleFunc("/photos", photos)
	http.HandleFunc("/mood", mood)
	http.HandleFunc("/projects", projects)
	http.HandleFunc("/about", about)
	http.HandleFunc("/wtf", wtf)
	http.HandleFunc("/health", health)

	log.Info("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Healthy as a horse!")
}
func wtf(w http.ResponseWriter, r *http.Request) {
	detect := mobiledetect.NewMobileDetect(r, nil)
	l := Layout{false, false, "", nil, detect.IsMobile()}
	utils.RenderTemplate(w, "wtf.html", l)
}
func projects(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "projects.html", nil)
}

func mood(w http.ResponseWriter, r *http.Request) {
	detect := mobiledetect.NewMobileDetect(r, nil)
	l := Layout{true, true, moodboardBackgroundGifs[rand.Intn(len(moodboardBackgroundGifs))], nil, detect.IsMobile()}
	if detect.IsMobile() {
		l.Photos = shufflePhotos(moodboardPagePhotos)[0:17]
	} else {
		l.Photos = shufflePhotos(moodboardPagePhotos)[0:29]
	}
	utils.RenderTemplate(w, "mood.html", l)
}

func about(w http.ResponseWriter, r *http.Request) {
	detect := mobiledetect.NewMobileDetect(r, nil)
	l := Layout{true, true, aboutBackgroundGifs[rand.Intn(len(aboutBackgroundGifs))], nil, detect.IsMobile()}
	utils.RenderTemplate(w, "about.html", l)
}

func home(w http.ResponseWriter, r *http.Request) {
	detect := mobiledetect.NewMobileDetect(r, nil)
	l := Layout{true, true, backgroundImages[rand.Intn(len(backgroundImages))], nil, detect.IsMobile()}
	utils.RenderTemplate(w, "home.html", l)
}

func photos(w http.ResponseWriter, r *http.Request) {
	detect := mobiledetect.NewMobileDetect(r, nil)
	l := Layout{false, false, "", nil, detect.IsMobile()}
	l.Photos = shufflePhotos(photosPagePhotos)
	utils.RenderTemplate(w, "photos.html", l)
}

func shufflePhotos(photos []Photo) []Photo {
	rand.Shuffle(len(photos), func(i, j int) { photos[i], photos[j] = photos[j], photos[i] })
	return photos
}

func refreshAboutGifs() {
	for {
		log.Info("Refreshing About Gifs")
		var s3Gifs = svc.RetrieveObjects("about/")
		var backgroundURLs []string
		for _, images := range s3Gifs {
			backgroundURLs = append(backgroundURLs, cloudFrontBaseURL+*images.Key)
		}
		aboutBackgroundGifs = backgroundURLs
		time.Sleep(5 * time.Minute)
	}
}

func refreshBackgroundImages() {
	for {
		log.Info("Refreshing Background Images")
		var s3Gifs = svc.RetrieveObjects("background/")
		var backgroundURLs []string
		for _, images := range s3Gifs {
			backgroundURLs = append(backgroundURLs, cloudFrontBaseURL+*images.Key)
		}
		backgroundImages = backgroundURLs
		time.Sleep(5 * time.Minute)
	}
}

func refreshMoodboardBackgroundGifs() {
	for {
		log.Info("Refreshing background gifs for moodboard")
		var s3Gifs = svc.RetrieveObjects("mood/gifs/")
		var gifURLs []string
		for _, gif := range s3Gifs {
			gifURLs = append(gifURLs, cloudFrontBaseURL+*gif.Key)
		}
		moodboardBackgroundGifs = gifURLs
		time.Sleep(5 * time.Minute)
	}
}

func refreshPhotosPagePhotos() {
	for {
		log.Info("Refreshing photos for Photos")
		s3Photos := svc.RetrieveObjects("website/")
		var photos []Photo
		for _, item := range s3Photos {
			id := strings.Replace(strings.Replace(*item.Key, ".jpg", "", 1), "website/", "", 1)
			p := Photo{cloudFrontBaseURL + *item.Key, "", id}
			photos = append(photos, p)
		}
		photosPagePhotos = photos
		time.Sleep(5 * time.Minute)
	}
}

func refreshMoodboardPagePhotos() {
	for {
		log.Info("Refreshing photos for Moodboard")
		s3Photos := svc.RetrieveObjects("mood/photos/")
		var photos []Photo
		for _, item := range s3Photos {
			id := strings.Replace(strings.Replace(*item.Key, ".jpg", "", 1), "mood/photos/", "", 1)
			captionChannel := make(chan string)
			go retrieveCaptionFromPhotoName(*item.Key, captionChannel)
			caption := <-captionChannel
			p := Photo{cloudFrontBaseURL + *item.Key, template.HTML(caption), id}
			photos = append(photos, p)
		}
		moodboardPagePhotos = photos
		time.Sleep(5 * time.Minute)
	}
}

func retrieveCaptionFromPhotoName(photoKey string, captionChannel chan string) {
	captionKey := strings.Replace(strings.Replace(photoKey, ".jpg", ".html", 1), "/photos/", "/captions/", 1)
	log.Debug("Retrieving Caption " + cloudFrontBaseURL + captionKey)
	resp, err := http.Get(cloudFrontBaseURL + captionKey)
	if err != nil || resp.StatusCode != http.StatusOK {
		captionChannel <- ""
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	responseString := string(bytes)
	captionChannel <- responseString
}

// Layout Object - drives html design
type Layout struct {
	DarkTheme     bool
	Background    bool
	BackgroundUrl string
	Photos        []Photo
	IsMobile      bool
}

type Photo struct {
	Url     string
	Caption template.HTML
	Id      string
}
