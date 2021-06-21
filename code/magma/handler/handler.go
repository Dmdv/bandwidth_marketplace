package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupHandlers sets up the necessary API end points.
func SetupHandlers(r *mux.Router) {
	r.HandleFunc("/", HomePageHandler)
	r.Handle("/metrics", promhttp.Handler())
}

// HomePageHandler provides basic info when accessing the home page of the server.
func HomePageHandler(w http.ResponseWriter, _ *http.Request) {
	page := "<div>I am a Magma Simulator </div>\n"

	if _, err := w.Write([]byte(page)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
