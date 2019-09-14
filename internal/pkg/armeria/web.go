package armeria

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func HandleScriptRead(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	ot := v["objectType"]
	on := v["objectName"]
	an := v["accessName"]
	ak := v["accessKey"]

	c := Armeria.characterManager.CharacterByName(an)
	if c == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if c.PasswordHash() != ak {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if ot != "mob" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := Armeria.mobManager.MobByName(on)
	if m == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s := ReadMobScript(m)

	_, _ = w.Write([]byte(s))
}

func HandleScriptWrite(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	ot := v["objectType"]
	on := v["objectName"]
	an := v["accessName"]
	ak := v["accessKey"]

	c := Armeria.characterManager.CharacterByName(an)
	if c == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if c.PasswordHash() != ak {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if ot != "mob" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := Armeria.mobManager.MobByName(on)
	if m == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	script, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	WriteMobScript(m, string(script))

	cp := c.Player()
	if cp != nil {
		cp.client.ShowColorizedText(
			fmt.Sprintf("The script has been saved to %s.", TextStyle(m.UnsafeName, TextStyleBold)),
			ColorSuccess,
		)
	}

	w.WriteHeader(http.StatusOK)
}

// InitWeb will initialize the HTTP web server, for serving the web client
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
		"/sfx/",
		"/gfx/",
		"/favicon.ico",
		"/scripteditor.html",
	}
	for _, route := range publicRoutes {
		r.PathPrefix(route).Handler(http.FileServer(http.Dir(Armeria.publicPath)))
	}
	r.PathPrefix("/oi/").Handler(http.StripPrefix("/oi/", http.FileServer(http.Dir(Armeria.objectImagesPath))))
	r.HandleFunc("/script/{objectType}/{objectName}/{accessName}/{accessKey}", HandleScriptRead).Methods("GET")
	r.HandleFunc("/script/{objectType}/{objectName}/{accessName}/{accessKey}", HandleScriptWrite).Methods("POST")
	r.PathPrefix("/ws").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(w, r)
	})
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, fmt.Sprintf("%s/index.html", Armeria.publicPath))
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port),
		handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(r),
	)

	if err != nil {
		Armeria.log.Fatal("error listening to http",
			zap.Error(err),
		)
	}
}
