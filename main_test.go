package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/obedtandadjaja/auth-go"
	"github.com/obedtandadjaja/auth-go/models/credential"
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

func TestCreateCredentialInvalidRequest(t *testing.T) {
	payload := []byte(`{"identifier":0,"password":0, "subject":0}`)

	req, _ := http.NewRequest("POST", "/credentials", bytes.NewBuffer(payload))
	rr := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, rr.Code)
}

func TestCreateCredential(t *testing.T) {
	clearCredentialsTable()

	rr := createCredential("email", "password", "website")

	checkResponseCode(t, http.StatusCreated, rr.Code)

	var responseBody map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &responseBody)

	if id, ok := responseBody["id"]; ok {
		credentials, _ := credential.All(app.DB)

		if len(credentials) == 1 {
			c := credentials[0]

			if c.Id != int(id.(float64)) {
				t.Errorf("Created credential id is wrong")
			}

			if c.Identifier != "email" {
				t.Errorf("Created credential identifier is wrong")
			}

			if c.Subject.String != "website" {
				t.Errorf("Created credential subject is wrong")
			}
		} else {
			t.Errorf("Expected one credential to be created, found %d", len(credentials))
		}
	} else {
		t.Errorf("Missing id in response")
	}
}

func createCredential(identifier, password, subject string) *httptest.ResponseRecorder {
	jsonString := fmt.Sprintf(`{"identifier":"%s","password":"%s","subject":"%s"}`, identifier, password, subject)

	payload := []byte(jsonString)

	req, _ := http.NewRequest("POST", "/credentials", bytes.NewBuffer(payload))
	rr := executeRequest(req)

	return rr
}

// Consider removing delete credential, since there is no use case for it. If we do want accounts
// to be deactivated, it should be a soft delete instead
func TestDeleteCredential(t *testing.T) {
	clearCredentialsTable()

	createCredential("email", "password", "website")

	rr := deleteCredential("email", "website")

	fmt.Println(rr)

	checkResponseCode(t, http.StatusNoContent, rr.Code)
}

func deleteCredential(identifier, subject string) *httptest.ResponseRecorder {
	jsonString := fmt.Sprintf(`{"identifier":"%s","subject":"%s"}`, identifier, subject)

	payload := []byte(jsonString)

	req, _ := http.NewRequest("DELETE", "/credentials", bytes.NewBuffer(payload))
	rr := executeRequest(req)

	return rr
}
