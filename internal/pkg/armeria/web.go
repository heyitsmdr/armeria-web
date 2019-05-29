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
	publicRoutes := []string{
		"/js/",
		"/css/",
		"/img/",
		"/vendor/",
		"/favicon.ico",
		"/scripteditor.html",
	}
	for _, route := range publicRoutes {
		r.PathPrefix(route).Handler(http.FileServer(http.Dir(Armeria.publicPath)))
	}
	r.PathPrefix("/oi/").Handler(http.StripPrefix("/oi/", http.FileServer(http.Dir(Armeria.objectImagesPath))))
	r.HandleFunc("/script/{objectType}/{objectName}/{accessName}/{accessKey}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("got " + vars["objectType"]))
		if err != nil {
			Armeria.log.Fatal("error writing to http for script",
				zap.Error(err),
			)
		}
	})
	r.PathPrefix("/ws").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(w, r)
	})
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, fmt.Sprintf("%s/index.html", Armeria.publicPath))
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		Armeria.log.Fatal("error listening to http",
			zap.Error(err),
		)
	}
}
