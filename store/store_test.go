package store

import (
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSaveURLMapping(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	const slug = "jCL9RStYD9"
	const longURL = "https://stripe.com"

	mock.ExpectExec("INSERT INTO url").WithArgs(slug, longURL).WillReturnResult(sqlmock.NewResult(1, 1))

	myDB := Database{sqlx.NewDb(db, "postgres")}
	err = myDB.SaveURLMapping(slug, longURL)
	if err != nil {
		t.Errorf("Unable to save URL mapping.")
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetLongURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	myDB := Database{sqlx.NewDb(db, "postgres")}

	rows := sqlmock.NewRows([]string{"destination"}).AddRow("https://stripe.com")
	mock.ExpectQuery("SELECT destination FROM url").WithArgs("jCL9RStYD9").WillReturnRows(rows)

	_, err = myDB.GetLongURL("jCL9RStYD9")
	if err != nil {
		t.Errorf("something went wrong: %s", err.Error())
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
