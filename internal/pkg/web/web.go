package web

import (
	"fmt"
	"log"
	"net/http"
)

// Init will initialize the HTTP web server
func Init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is Armeria.")
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
