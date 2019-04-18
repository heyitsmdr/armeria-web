package armeria

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Init will initialize the HTTP web server, for serving the web client
func InitWeb(state *GameState, publicPath string) {
	log.Printf("[web] serving client from: %s", publicPath)

	r := mux.NewRouter()

	// Set up routes
	r.PathPrefix("/js").Handler(http.FileServer(http.Dir(publicPath)))
	r.PathPrefix("/css").Handler(http.FileServer(http.Dir(publicPath)))
	r.PathPrefix("/img").Handler(http.FileServer(http.Dir(publicPath)))
	r.PathPrefix("/favicon.ico").Handler(http.FileServer(http.Dir(publicPath)))
	r.PathPrefix("/ws").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(state, w, r)
	})
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", publicPath))
	})

	http.Handle("/", r)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("[web] ListenAndServe: ", err)
	}
}
