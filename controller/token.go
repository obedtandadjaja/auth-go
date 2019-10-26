package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/obedtandadjaja/auth-go/auth/jwt"
	"github.com/obedtandadjaja/auth-go/models/credential"
)

const (
	MAX_FAILED_ATTEMPTS = 3
)

type TokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	Jwt string `json:"jwt"`
}

func Token(sr *SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseTokenRequest(r)
	if err != nil {
		return HandlerError{400, err, nil}
	}

	response, err := processTokenRequest(sr, request)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseTokenRequest(r *http.Request) (*TokenRequest, error) {
	var request TokenRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processTokenRequest(sr *SharedResources, request *TokenRequest) (*TokenResponse, error) {
	var response TokenResponse

	credential, err := credential.FindBy(sr.DB, map[string]interface{}{
		"id": request.RefreshToken,
	})
	if err != nil {
		return &response, HandlerError{401, errors.New("Invalid refresh token"), err}
	}

	tokenString, err := jwt.Generate(credential.Id)
	if err != nil {
		return &response, HandlerError{500, err, err}
	}

	response.Jwt = tokenString
	return &response, nil
}
