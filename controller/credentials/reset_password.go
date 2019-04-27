package credentials

import (
	"net/http"
	"encoding/json"
	"errors"
	"database/sql"

	"github.com/obedtandadjaja/auth-go/controller"
	"github.com/obedtandadjaja/auth-go/models/credential"
)

type ResetPasswordRequest struct {
	Identifier         string `json:"identifier"`
	Subject            string `json:"subject"`
	PasswordResetToken string `json:"password_reset_token"`
	NewPassword        string `json:"new_password"`
}

type ResetPasswordResponse struct {
	Id int `json:"id"`
}

func ResetPasssword(sr *controller.SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseResetPassword(r)
	if err != nil {
		return controller.HandlerError{400, err}
	}

	response, err := processResetPasswordRequest(sr, request, r)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseResetPassword(r * http.Request) (*ResetPasswordRequest, error) {
	var request ResetPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processResetPasswordRequest(sr *controller.SharedResources, request *ResetPasswordRequest, r *http.Request) (*ResetPasswordResponse, error) {
	var response ResetPasswordResponse

	cred, err := credential.FindBy(sr.DB, map[string]interface{}{
		"identifier": request.Identifier,
		"subject": request.Subject,
	})
	if err != nil {
		return &response, controller.HandlerError{404, errors.New("Credential not found")}
	}

	if !cred.PasswordResetToken.Valid {
		return &response, controller.HandlerError{400, errors.New("Credential did not apply for password reset")}
	}

	if cred.PasswordResetToken.String != request.PasswordResetToken {
		return &response, controller.HandlerError{401, errors.New("Wrong password reset token")}
	}

	cred.Password = sql.NullString{String: request.NewPassword, Valid: true}
	return &response, nil
}
