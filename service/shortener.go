// Package service provides the business logic for URL shortener
package service

import (
	"time"

	"github.com/khachikyan/shortify/base62"
	"github.com/khachikyan/shortify/entity"
	"github.com/khachikyan/shortify/store"

	"github.com/sony/sonyflake"
)

type Shortener interface {
	Shorten(longURL string) (string, error)
	GetLongURL(slug string) (string, error)
	GetURLVisits(slug string, unit string) (entity.Metric, error)
	LogURLVisit(slug string, date string) error
}

type Store struct {
	storage store.StoreInterface
}

func New(s store.StoreInterface) *Store {
	return &Store{s}
}

// Generates an unique slug for a URL
func generateSlug() (string, error) {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		return "", err
	}

	slug := base62.Encode(id)
	return slug, nil
}

// Shorten creates a slug for the short URL and persists it in the store
func (s *Store) Shorten(longURL string) (string, error) {
	slug, err := generateSlug()
	if err != nil {
		return "", err
	}
	err = s.storage.SaveURLMapping(slug, longURL)
	return slug, err
}

// GetLongURL retrieves the long URL from the store
func (s *Store) GetLongURL(slug string) (string, error) {
	longURL, err := s.storage.GetLongURL(slug)
	return longURL, err
}

// GetURLVisits quries the store for URL visits
func (s *Store) GetURLVisits(slug string, unit string) (entity.Metric, error) {
	var queryDate string
	currentDate := time.Now()

	if unit == "day" {
		queryDate = currentDate.AddDate(0, 0, -1).Format("2006-01-02 15:04:05")
	} else if unit == "week" {
		queryDate = currentDate.AddDate(0, 0, -7).Format("2006-01-02 15:04:05")
	} else {
		// HACK
		queryDate = time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC).Format("2006-01-02 15:04:05")
	}

	return s.storage.GetURLMetrics(slug, queryDate)
}

// LogURLVisit inserts URL visit event in the store
func (s *Store) LogURLVisit(slug string, date string) error {
	return s.storage.SaveURLVist(slug, date)
}
