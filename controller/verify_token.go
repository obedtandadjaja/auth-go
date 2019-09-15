package controller

import (
	"encoding/json"
	"net/http"

	"github.com/obedtandadjaja/auth-go/auth/jwt"
)

type VerifyRequest struct {
	Jwt string `json:"jwt"`
}

type VerifyResponse struct {
	Verified bool `json:"verified"`
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

	err := jwt.Verify(request.Jwt)

	response.Verified = err != nil
	return &response, nil
}
