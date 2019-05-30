package credentials

import (
	"net/http"
	"encoding/json"
	"errors"

	"github.com/obedtandadjaja/auth-go/controller"
	"github.com/obedtandadjaja/auth-go/models/credential"
)

type InitiatePasswordResetRequest struct {
	Identifier string `json:"identifier"`
	Subject    string `json:"subject"`
}

type InitiatePasswordResetResponse struct {
	PasswordResetToken string `json:"password_reset_token"`
}

func InitiatePasswordReset(sr *controller.SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseInitiatePasswordResetRequest(r)
	if err != nil {
		return controller.HandlerError{400, err, err}
	}

	response, err := processInitiatePasswordResetRequest(sr, request, r)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseInitiatePasswordResetRequest(r *http.Request) (*InitiatePasswordResetRequest, error) {
	var request InitiatePasswordResetRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processInitiatePasswordResetRequest(sr *controller.SharedResources, request *InitiatePasswordResetRequest, r *http.Request) (*InitiatePasswordResetResponse, error) {
	var response InitiatePasswordResetResponse

	cred, err := credential.FindBy(sr.DB, map[string]interface{}{
		"identifier": request.Identifier,
		"subject": request.Subject,
	})
	if err != nil {
		return &response, controller.HandlerError{404, errors.New("Credential not found"), err}
	}

	err = cred.SetPasswordResetToken(sr.DB)
	if err != nil {
		return &response, controller.HandlerError{400, errors.New("Failed to initiate password reset"), err}
	}

	response.PasswordResetToken = cred.PasswordResetToken.String
	return &response, nil
}
