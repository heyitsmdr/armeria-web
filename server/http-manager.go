package arcadia

import (
	"fmt"
	"log"
	"net/http"
)

// HTTPManager oversees all http connections to the server
type HTTPManager struct {
}

// NewHTTPManager creates a new HttpManager instance
func NewHTTPManager() *HTTPManager {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is Arcadia.")
	})

	return &HTTPManager{}
}

func (*HTTPManager) ReadyToServe() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
