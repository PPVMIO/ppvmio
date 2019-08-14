package handler


import (
	"html/template"
	"path/filepath"
	"log"
	"net/http"
)
	
type Photo struct {
	Path	string
}

type Todo struct {
	Name        string
	Description string
}

func HandlePhotos(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("web", "layout.html")
	
	log.Println(r.URL.Path)
	fp := filepath.Join("web", "photos/index.html")
	
	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	photo := Photo{"hello"}
	
	if err := tmpl.ExecuteTemplate(w, "layout", photo); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
