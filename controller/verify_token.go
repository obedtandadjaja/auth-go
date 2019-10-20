package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/obedtandadjaja/auth-go/auth/jwt"
)

type VerifyRequest struct {
	Jwt string `json:"jwt"`
}

type VerifyResponse struct {
	CredentialId string `json:"credential_id"`
	Verified     bool   `json:"verified"`
}

func Verify(sr *SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseVerifyRequest(r)
	if err != nil {
		return HandlerError{400, err, nil}
	}

	response, err := processVerifyRequest(sr, request)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseVerifyRequest(r *http.Request) (*VerifyRequest, error) {
	var request VerifyRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processVerifyRequest(sr *SharedResources, request *VerifyRequest) (*VerifyResponse, error) {
	var response VerifyResponse

	credentialId, err := jwt.Verify(request.Jwt)

	if err != nil {
		return &response, HandlerError{400, errors.New("Invalid JWT token"), nil}
	} else {
		response.CredentialId = credentialId
		response.Verified = true
	}
	return &response, nil
}
