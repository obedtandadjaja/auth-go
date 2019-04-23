package controller

import (
	"net/http"
	"encoding/json"
	"errors"

	"github.com/obedtandadjaja/auth-go/models/credential"
	"github.com/obedtandadjaja/auth-go/auth/jwt"
	"github.com/obedtandadjaja/auth-go/auth/hash"
)

type TokenRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type TokenResponse struct {
	Jwt string `json:"jwt"`
}

func Token(sr *SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseRequest(r)
	if err != nil {
		return HandlerError{400, err}
	}

	response, err := processRequest(sr, request)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseRequest(r *http.Request) (*TokenRequest, error) {
	var request TokenRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processRequest(sr *SharedResources, request *TokenRequest) (*TokenResponse, error) {
	var response TokenResponse

	credential, err := credential.FindBy(sr.DB, "identifier", request.Identifier)
	if err != nil {
		return &response, HandlerError{401, errors.New("Invalid credentials")}
	}

	if hashValue := credential.Password.String; !hash.ValidatePasswordHash(request.Password, hashValue) {
		return &response, HandlerError{401, errors.New("Invalid credentials")}
	}

	tokenString, err := jwt.Generate(request.Identifier)
	if err != nil {
		return &response, HandlerError{500, err}
	}

	response.Jwt = tokenString
	return &response, nil
}
