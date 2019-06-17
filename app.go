package main

import (
	"fmt"
	"io"
	"net/http"
	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	
	homeDir := http.Dir("./web/home")
	homeHandler := http.StripPrefix("/home", http.FileServer(homeDir))

	infraDir := http.Dir("./web/infrastructure")
	infraHandler := http.StripPrefix("/infrastructure", http.FileServer(infraDir))

	assetsDir := http.Dir("./assets")
	assetsHandler := http.StripPrefix("/assets", http.FileServer(assetsDir))

	r.PathPrefix("/home").Handler(homeHandler).Methods("GET")
	r.PathPrefix("/assets").Handler(assetsHandler).Methods("GET")
	r.PathPrefix("/infrastructure").Handler(infraHandler).Methods("GET")
	r.HandleFunc("/health", HealthCheckHandler)

	return r
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    io.WriteString(w, `{"alive": true}`)
}


func main() {
	r := newRouter()
	http.ListenAndServe(":8080", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}