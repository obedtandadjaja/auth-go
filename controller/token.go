package controller

import (
	"encoding/json"
	"net/http"

	"github.com/obedtandadjaja/auth-go/auth/jwt"
	"github.com/obedtandadjaja/auth-go/models/credential"
	"github.com/obedtandadjaja/auth-go/models/session"
)

const (
	MAX_FAILED_ATTEMPTS = 3
)

type TokenRequest struct {
	SessionJwt string `json:"session"`
}

type TokenResponse struct {
	Jwt string `json:"jwt"`
}

func Token(sr *SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseTokenRequest(r)
	if err != nil {
		return HandlerError{400, "", err}
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

	refreshTokenUuid, err := jwt.VerifyRefreshToken(request.SessionJwt)
	if err != nil {
		return &response, HandlerError{401, "Invalid refresh token", err}
	}

	// find the refresh token record
	sessionRecord, err := session.FindBy(sr.DB, map[string]interface{}{
		"uuid": refreshTokenUuid,
	})
	if err != nil {
		return &response, HandlerError{401, "Invalid refresh token", err}
	}

	// find the credential record
	credential, err := credential.FindBy(sr.DB, map[string]interface{}{
		"id": sessionRecord.CredentialId,
	})
	if err != nil {
		return &response, HandlerError{404, "Invalid refresh token", err}
	}

	tokenString, err := jwt.GenerateAccessToken(credential.Uuid)
	if err != nil {
		return &response, HandlerError{500, "Internal Server Error", err}
	}

	response.Jwt = tokenString
	return &response, nil
}
