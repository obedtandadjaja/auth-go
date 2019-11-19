package credentials

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/obedtandadjaja/auth-go/auth/jwt"
	"github.com/obedtandadjaja/auth-go/controller"
	"github.com/obedtandadjaja/auth-go/models/credential"
	"github.com/obedtandadjaja/auth-go/models/session"
)

type CreateRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type CreateResponse struct {
	CredentialUuid string `json:"credential_uuid"`
	Jwt            string `json:"jwt"`
	SessionJwt     string `json:"session"`
}

func Create(sr *controller.SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseCreateRequest(r)
	if err != nil {
		return controller.HandlerError{400, err.Error(), err}
	}

	response, err := processCreateRequest(sr, request, r)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseCreateRequest(r *http.Request) (*CreateRequest, error) {
	var request CreateRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if request.Password == "" {
		return &request, errors.New("Missing required field")
	}

	return &request, err
}

func processCreateRequest(sr *controller.SharedResources, request *CreateRequest, r *http.Request) (*CreateResponse, error) {
	var response CreateResponse

	cred := credential.Credential{
		Email:    sql.NullString{String: request.Email, Valid: request.Email != ""},
		Phone:    sql.NullString{String: request.Phone, Valid: request.Phone != ""},
		Password: sql.NullString{String: request.Password, Valid: true},
	}

	err := cred.Create(sr.DB)
	if err != nil {
		return &response, controller.HandlerError{
			400,
			"Failed to create credential",
			err,
		}
	}

	newSession := session.Session{
		CredentialId: cred.Id,
		ExpiresAt:    time.Now().Add(time.Duration(24 * 180 * time.Hour)),
		IpAddress:    sql.NullString{String: r.RemoteAddr, Valid: true},
		UserAgent:    sql.NullString{String: r.UserAgent(), Valid: true},
	}
	err = newSession.Create(sr.DB)
	if err != nil {
		return &response, controller.HandlerError{500, "Internal Server Error", err}
	}

	refreshTokenChan := make(chan string)
	accessTokenChan := make(chan string)

	go func() {
		refreshTokenJwt, _ := jwt.GenerateRefreshToken(newSession.Uuid)
		refreshTokenChan <- refreshTokenJwt
	}()

	go func() {
		accessTokenJwt, _ := jwt.GenerateAccessToken(cred.Uuid)
		accessTokenChan <- accessTokenJwt
	}()

	response.Jwt = <-accessTokenChan
	response.SessionJwt = <-refreshTokenChan
	response.CredentialUuid = cred.Uuid
	return &response, nil
}
