package credentials

import (
	"net/http"
	"encoding/json"
	// "errors"
	"database/sql"

	"github.com/obedtandadjaja/auth-go/controller"
	"github.com/obedtandadjaja/auth-go/models/credential"
)

type CreateRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
	Subject    string `json:"subject"`
}

type CreateResponse struct {
	Id int64 `json:"id"`
}

func Create(sr *controller.SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseRequest(r)
	if err != nil {
		return controller.HandlerError{400, err}
	}

	response, err := processRequest(sr, request)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseRequest(r *http.Request) (*CreateRequest, error) {
	var request CreateRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processRequest(sr *controller.SharedResources, request *CreateRequest) (*CreateResponse, error) {
	var response CreateResponse

	cred := credential.Credential{
		Identifier: request.Identifier,
		Password: sql.NullString{String: request.Password, Valid: true},
		Subject: sql.NullString{String: request.Subject, Valid: true},
	}

	err := cred.Create(sr.DB)
	if err != nil {
		// return &response, controller.HandlerError{400, errors.New("Failed to create credential")}
		return &response, controller.HandlerError{400, err}
	}

	response.Id = cred.Id
	return &response, nil
}
