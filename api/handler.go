package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/khachikyan/shortify/service"
)

type shortenReq struct {
	LongURL string `json:"long_url"`
}

type shortenRes struct {
	ShortURL string `json:"short_url"`
}

type Api struct {
	service.Shortener
}

func New(a service.Shortener) *Api {
	return &Api{a}
}

// ShortURL Route handler for encoding a long url
func (api *Api) ShortURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var payload shortenReq
	err := decoder.Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	slug, err := api.Shortener.Shorten(payload.LongURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	res := shortenRes{
		ShortURL: "localhost:80/" + slug,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusOK)
}

// Redirect Route handler for redirecting short url to long url
func (api *Api) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	longURL, err := api.GetLongURL(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	currentDate := time.Now().Format("2006-01-02 15:04:05")
	api.LogURLVisit(slug, currentDate)

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}

// URLMetrics Route handle for short URL analytics
func (api *Api) URLMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	unit := r.URL.Query().Get("unit")

	// Set default query parameter to 24 hours
	if unit == "" {
		unit = "day"
	}

	count, _ := api.GetURLVisits(slug, unit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(count)
	w.WriteHeader(http.StatusOK)
}
