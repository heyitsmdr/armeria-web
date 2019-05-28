package armeria

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

func ConfigureStaticRoute(r *mux.Router, pathPrefix string, dir string) {
	s := http.StripPrefix(pathPrefix, http.FileServer(http.Dir(dir+"/")))
	r.PathPrefix(pathPrefix).Handler(s)
}

// Init will initialize the HTTP web server, for serving the web client
func InitWeb(port int) {
	Armeria.log.Info("serving http requests",
		zap.String("path", Armeria.publicPath),
		zap.Int("port", port),
	)

	// Set up routes
	r := mux.NewRouter()
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir(Armeria.publicPath)))
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir(Armeria.publicPath)))
	r.PathPrefix("/img/").Handler(http.FileServer(http.Dir(Armeria.publicPath)))
	r.PathPrefix("/oi/").Handler(http.StripPrefix("/oi/", http.FileServer(http.Dir(Armeria.objectImagesPath))))
	r.PathPrefix("/ws").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(w, r)
	})
	r.PathPrefix("/favicon.ico").Handler(http.FileServer(http.Dir(Armeria.publicPath)))
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", Armeria.publicPath))
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		Armeria.log.Fatal("error listening to http",
			zap.Error(err),
		)
	}
}
