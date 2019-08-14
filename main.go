package main

import (
	"log"
	"net/http"
	"ppvmio/utils"
	"ppvmio/svc"
	"time"
	"math/rand"
)

func main() {
	utils.LoadTemplates()

	http.HandleFunc("/home", home)
	// http.HandleFunc("/infrastructure", infrastructure)
	http.HandleFunc("/photos", photos)
	
	log.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	err := utils.RenderTemplate(w, "home.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func photos(w http.ResponseWriter, r *http.Request) {

	s3Photos := svc.RetrievePhotos(5)
	log.Println(len(s3Photos))
	var photos []Photo

	for _, item := range s3Photos {
		p := Photo{"http://photos-ppvmio-public.s3.amazonaws.com/" + *item.Key}
		log.Println("Adding Photo: " + p.Path)
		photos = append(photos, p)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(photos), func(i, j int) { photos[i], photos[j] = photos[j], photos[i] })

	err := utils.RenderTemplate(w, "photos.html", photos)
	if err != nil {
		log.Println(err)
	}
}


type Photo struct {
	Path	string
}