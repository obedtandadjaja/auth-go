package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/obedtandadjaja/auth-go/auth/hash"
	"github.com/obedtandadjaja/auth-go/auth/jwt"
	"github.com/obedtandadjaja/auth-go/models/credential"
	"github.com/obedtandadjaja/auth-go/models/session"
)

type LoginRequest struct {
	CredentialUuid string `json:"credential_uuid"`
	Password       string `json:"password"`
}

type LoginResponse struct {
	Jwt        string `json:"jwt"`
	SessionJwt string `json:"session"`
}

func Login(sr *SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseLoginRequest(r)
	if err != nil {
		return HandlerError{400, "", err}
	}

	response, err := processLoginRequest(sr, request, r)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseLoginRequest(r *http.Request) (*LoginRequest, error) {
	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processLoginRequest(sr *SharedResources, request *LoginRequest, r *http.Request) (*LoginResponse, error) {
	var response LoginResponse

	credential, err := credential.FindBy(sr.DB, map[string]interface{}{
		"uuid": request.CredentialUuid,
	})
	if err != nil {
		return &response, HandlerError{401, "Invalid credentials", err}
	}

	if credential.LockedUntil.Valid && credential.LockedUntil.Time.After(time.Now()) {
		return &response, HandlerError{
			401,
			fmt.Sprintf("Locked until %v", credential.LockedUntil.Time.Sub(time.Now())),
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

		return &response, HandlerError{401, "Invalid credentials", nil}
	}

	// TODO: move this as a goroutine
	// don't care about this error if there is any
	err = credential.Update(sr.DB, map[string]interface{}{
		"failed_attempts": 0,
		"locked_until":    nil,
	})

	newSession := session.Session{
		CredentialId: credential.Id,
		ExpiresAt:    time.Now().Add(time.Duration(24 * 180 * time.Hour)),
		IpAddress:    sql.NullString{String: r.RemoteAddr, Valid: true},
		UserAgent:    sql.NullString{String: r.UserAgent(), Valid: true},
	}
	err = newSession.Create(sr.DB)
	if err != nil {
		return &response, HandlerError{500, "Internal Server Error", err}
	}

	sessionJwt, err := jwt.GenerateRefreshToken(newSession.Uuid)
	if err != nil {
		return &response, HandlerError{500, "Internal Server Error", err}
	}

	accessTokenJwt, err := jwt.GenerateAccessToken(credential.Uuid)
	if err != nil {
		return &response, HandlerError{500, "Internal Server Error", err}
	}

	response.Jwt = accessTokenJwt
	response.SessionJwt = sessionJwt
	return &response, nil
}
