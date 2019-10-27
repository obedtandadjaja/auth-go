package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/obedtandadjaja/auth-go/auth/hash"
	"github.com/obedtandadjaja/auth-go/auth/jwt"
	"github.com/obedtandadjaja/auth-go/models/credential"
	"github.com/obedtandadjaja/auth-go/models/refresh_token"
)

type RefreshTokenRequest struct {
	CredentialUuid string `json:"credential_uuid"`
	Password       string `json:"password"`
}

type RefreshTokenResponse struct {
	Jwt          string `json:"jwt"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func RefreshToken(sr *SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseRefreshTokenRequest(r)
	if err != nil {
		return HandlerError{400, err, nil}
	}

	response, err := processRefreshTokenRequest(sr, request, r)
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

func processRefreshTokenRequest(sr *SharedResources, request *RefreshTokenRequest, r *http.Request) (*RefreshTokenResponse, error) {
	var response RefreshTokenResponse

	credential, err := credential.FindBy(sr.DB, map[string]interface{}{
		"uuid": request.CredentialUuid,
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

	refreshToken := refresh_token.RefreshToken{
		CredentialId: credential.Id,
		ExpiresAt:    time.Now().Add(time.Duration(24 * 180 * time.Hour)),
		IpAddress:    sql.NullString{String: r.RemoteAddr, Valid: true},
		UserAgent:    sql.NullString{String: r.UserAgent(), Valid: true},
	}
	err = refreshToken.Create(sr.DB)
	if err != nil {
		return &response, HandlerError{500, errors.New("Internal Server Error"), err}
	}

	refreshTokenJwt, err := jwt.GenerateRefreshToken(refreshToken.Uuid)
	if err != nil {
		return &response, HandlerError{500, err, err}
	}

	accessTokenJwt, err := jwt.GenerateAccessToken(credential.Uuid)
	if err != nil {
		return &response, HandlerError{500, err, err}
	}

	response.Jwt = accessTokenJwt
	response.RefreshToken = refreshTokenJwt
	return &response, nil
}
