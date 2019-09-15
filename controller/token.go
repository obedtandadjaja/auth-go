package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/obedtandadjaja/auth-go/auth/hash"
	"github.com/obedtandadjaja/auth-go/auth/jwt"
	"github.com/obedtandadjaja/auth-go/models/credential"
)

const (
	MAX_FAILED_ATTEMPTS = 3
)

type TokenRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
	Subject    string `json:"subject"`
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
		"identifier": request.Identifier,
		"subject":    request.Subject,
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

	// don't care about this error if there is any
	err = credential.Update(sr.DB, map[string]interface{}{
		"failed_attempts": 0,
		"locked_until":    nil,
	})

	tokenString, err := jwt.Generate(request.Identifier)
	if err != nil {
		return &response, HandlerError{500, err, err}
	}

	response.Jwt = tokenString
	return &response, nil
}
