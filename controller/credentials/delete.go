package credentials

import (
	"net/http"
	"encoding/json"
	"errors"

	"github.com/obedtandadjaja/auth-go/controller"
	"github.com/obedtandadjaja/auth-go/models/credential"
)

type DeleteRequest struct {
	Identifier string `json:"identifier"`
	Subject    string `json:"subject"`
}

type DeleteResponse struct {
	Id int `json:"id"`
}

func Delete(sr *controller.SharedResources, w http.ResponseWriter, r * http.Request) error {
	request, err := parseDeleteRequest(r)
	if err != nil {
		return controller.HandlerError{400, err}
	}

	response, err := processDeleteRequest(sr, request,r)
	if err != nil {
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseDeleteRequest(r *http.Request) (*DeleteRequest, error) {
	var request DeleteRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processDeleteRequest(sr *controller.SharedResources, request *DeleteRequest, r *http.Request) (*DeleteResponse, error) {
	var response DeleteResponse

	cred, err := credential.FindBy(sr.DB, "identifier", request.Identifier)
	if err != nil {
		return &response, controller.HandlerError{404, errors.New("Credential not found")}
	}

	err = cred.Delete(sr.DB)
	if err != nil {
		return &response, controller.HandlerError{400, errors.New("Failed to delete credential")}
	}

	response.Id = cred.Id
	return &response, nil
}
