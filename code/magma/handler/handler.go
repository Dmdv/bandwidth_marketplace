package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

// SetupHandlers sets up the necessary API end points.
func SetupHandlers(r *mux.Router) {
	r.HandleFunc("/", HomePageHandler)
}

// HomePageHandler provides basic info when accessing the home page of the server.
func HomePageHandler(w http.ResponseWriter, _ *http.Request) {
	page := "<div>I am a Magma Simulator </div>\n"

	if _, err := w.Write([]byte(page)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
