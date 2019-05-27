package armeria

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

// Init will initialize the HTTP web server, for serving the web client
func InitWeb(port int) {
	Armeria.log.Info("serving http requests",
		zap.String("path", Armeria.publicPath),
		zap.Int("port", port),
	)
	r := mux.NewRouter()

	// Set up routes
	r.PathPrefix("/js").Handler(http.FileServer(http.Dir(Armeria.publicPath)))
	r.PathPrefix("/css").Handler(http.FileServer(http.Dir(Armeria.publicPath)))
	r.PathPrefix("/img").Handler(http.FileServer(http.Dir(Armeria.publicPath)))
	r.PathPrefix("/favicon.ico").Handler(http.FileServer(http.Dir(Armeria.publicPath)))
	r.PathPrefix("/ws").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(w, r)
	})
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", Armeria.publicPath))
	})

	http.Handle("/", r)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		Armeria.log.Fatal("error listening to http",
			zap.Error(err),
		)
	}
}
