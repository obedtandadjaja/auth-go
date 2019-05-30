package main_test

import (
	"os"
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/obedtandadjaja/auth-go"
)

var app main.App

func TestMain(m *testing.M) {
	app = main.App{}
	app.Initialize(
		os.Getenv("ENV"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
	)

	code := m.Run()

	os.Exit(code)
}

func clearCredentialsTable() {
	app.DB.Exec("delete from credentials")
	app.DB.Exec("alter sequence credentials_id_seq restart with 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestCreateCredential(t *testing.T) {
	clearCredentialsTable()

	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/credentials", nil)

	app.Router.ServeHTTP(rr, req)

	checkResponseCode(t, http.StatusBadRequest, rr.Code)
}
