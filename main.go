package main

import (
	"net/http"

	"github.com/khachikyan/shortify/api"
	"github.com/khachikyan/shortify/entity"
	"github.com/khachikyan/shortify/service"
	"github.com/khachikyan/shortify/store"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

const (
	configFilePath = "./config.toml"
)

// Decodes TOML configuration file, and returns a config object
func parseTOML() (*entity.Config, error) {
	var config entity.Config
	_, err := toml.DecodeFile(configFilePath, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	config, err := parseTOML()
	if err != nil {
		return
	}
	dbConn, err := store.New(config)
	if err != nil {
		return
	}

	URLService := service.New(dbConn)
	URLApi := api.New(URLService)

	r := mux.NewRouter()

	// API Endpoints
	r.HandleFunc("/short", URLApi.ShortURL).Methods(http.MethodPost)
	r.HandleFunc("/{slug}", URLApi.Redirect).Methods(http.MethodGet)
	r.HandleFunc("/{slug}/clicks", URLApi.URLMetrics).Methods(http.MethodGet)

	// Initialize HTTP middleware
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	n.UseHandler(r)

	http.ListenAndServe(":8000", n)

}
