package credentials

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/obedtandadjaja/auth-go/controller"
	"github.com/obedtandadjaja/auth-go/models/credential"

	"github.com/lib/pq"
)

type CreateRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
	Subject    string `json:"subject"`
}

type CreateResponse struct {
	Id int `json:"id"`
}

func Create(sr *controller.SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseCreateRequest(r)
	if err != nil {
		return controller.HandlerError{400, err, err}
	}

	response, err := processCreateRequest(sr, request, r)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseCreateRequest(r *http.Request) (*CreateRequest, error) {
	var request CreateRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processCreateRequest(sr *controller.SharedResources, request *CreateRequest, r *http.Request) (*CreateResponse, error) {
	var response CreateResponse

	cred := credential.Credential{
		Identifier: request.Identifier,
		Password:   sql.NullString{String: request.Password, Valid: true},
		Subject:    sql.NullString{String: request.Subject, Valid: true},
		IpAddress:  sql.NullString{String: r.RemoteAddr, Valid: true},
	}

	err := cred.Create(sr.DB)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return &response, controller.HandlerError{
				400,
				errors.New("There is already an existing credential with this identifier"),
				err,
			}
		} else {
			return &response, controller.HandlerError{
				400,
				errors.New("Failed to create credential"),
				err,
			}
		}
	}

	response.Id = cred.Id
	return &response, nil
}
