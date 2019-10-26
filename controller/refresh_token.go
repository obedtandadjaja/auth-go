package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/obedtandadjaja/auth-go/auth/hash"
	"github.com/obedtandadjaja/auth-go/models/credential"
)

type RefreshTokenRequest struct {
	CredentialId string `json:"credential_id"`
	Password     string `json:"password"`
}

type RefreshTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
}

func RefreshToken(sr *SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseRefreshTokenRequest(r)
	if err != nil {
		return HandlerError{400, err, nil}
	}

	response, err := processRefreshTokenRequest(sr, request)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseRefreshTokenRequest(r *http.Request) (*RefreshTokenRequest, error) {
	var request RefreshTokenRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processRefreshTokenRequest(sr *SharedResources, request *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	var response RefreshTokenResponse

	credential, err := credential.FindBy(sr.DB, map[string]interface{}{
		"id": request.CredentialId,
	})
	if err != nil {
		return &response, HandlerError{401, errors.New("Invalid credentials"), err}
	}

	if credential.LockedUntil.Valid && credential.LockedUntil.Time.After(time.Now()) {
		return &response, HandlerError{
			401,
			errors.New(fmt.Sprintf("Locked until %v", credential.LockedUntil.Time.Sub(time.Now()))),
			nil,
		}
	}

	if hashValue := credential.Password.String; !hash.ValidatePasswordHash(request.Password, hashValue) {
		if credential.FailedAttempts == MAX_FAILED_ATTEMPTS {
			credential.Update(sr.DB, map[string]interface{}{
				"locked_until": time.Now().Add(time.Duration(credential.FailedAttempts*10) * time.Minute),
			})
		}
		credential.IncrementFailedAttempt(sr.DB)

		return &response, HandlerError{401, errors.New("Invalid credentials"), nil}
	}

	// TODO: move this as a goroutine
	// don't care about this error if there is any
	err = credential.Update(sr.DB, map[string]interface{}{
		"failed_attempts": 0,
		"locked_until":    nil,
	})

	return &response, nil
}
