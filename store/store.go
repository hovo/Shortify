package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/khachikyan/shortify/entity"
	_ "github.com/lib/pq"
)

type StoreInterface interface {
	GetLongURL(slug string) (string, error)
	SaveURLMapping(slug string, longURL string) error
	GetURLMetrics(slug string, date string) (entity.Metric, error)
	SaveURLVist(slug string, date string) error
}

type Database struct {
	*sqlx.DB
}

// New initializes a new database connection
func New(config *entity.Config) (*Database, error) {
	dbConnInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host, config.DB.Port, config.DB.Username, config.DB.Password, config.DB.Database)

	var db *sqlx.DB
	db, err := sqlx.Open("postgres", dbConnInfo)
	if err != nil {
		panic(err)
	}

	tx := db.MustBegin()
	query := "CREATE TABLE IF NOT EXISTS url (id serial, slug char(10), destination varchar(255));"
	tx.MustExec(query)

	query = "CREATE TABLE IF NOT EXISTS visits (id serial, slug char(10), visited_at timestamp);"
	tx.MustExec(query)

	err = tx.Commit()

	return &Database{db}, err
}

// GetLongURL retrieves the destination/long url corresponding to short url slug
func (db *Database) GetLongURL(slug string) (string, error) {
	var destinationURL string

	query := "SELECT destination FROM url WHERE slug=$1"
	row := db.QueryRowx(query, slug)
	err := row.Scan(&destinationURL)
	if err != nil {
		return "", err
	}

	return destinationURL, nil
}

// SaveURLMapping insert new short url slug to long url entry
func (db *Database) SaveURLMapping(slug string, longURL string) error {
	query := "INSERT INTO url (slug, destination) VALUES ($1, $2);"
	_, err := db.Exec(query, slug, longURL)
	return err
}

// GetURLMetrics retrieves the number of times a short url has been accessed
func (db *Database) GetURLMetrics(slug string, date string) (entity.Metric, error) {
	var URLCount entity.Metric
	query := "SELECT COUNT(id) FROM visits where slug=$1 and visited_at >= $2;"

	row := db.QueryRowx(query, slug, date)
	err := row.Scan(&URLCount.Visits)

	return URLCount, err
}

// SaveURLVist inserts URL visit event entry in the store
func (db *Database) SaveURLVist(slug string, date string) error {
	query := "INSERT INTO visits (slug, visited_at) VALUES ($1, $2);"
	_, err := db.Exec(query, slug, date)
	return err
}
