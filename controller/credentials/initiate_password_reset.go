package credentials

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/obedtandadjaja/auth-go/controller"
	"github.com/obedtandadjaja/auth-go/models/credential"
)

type InitiatePasswordResetRequest struct {
	CredentialUuid string `json:"credential_uuid"`
}

type InitiatePasswordResetResponse struct {
    ResetPasswordToken string `json:"reset_password_token"`
}

func InitiatePasswordReset(sr *controller.SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseInitiatePasswordResetRequest(r)
	if err != nil {
		return controller.HandlerError{400, "", err}
	}

    response, err := processInitiatePasswordResetRequest(sr, request, r)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
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
		"id": request.CredentialUuid,
	})
	if err != nil {
		return &response, controller.HandlerError{404, "Credential not found", err}
	}

	// token is the last 6 digit of the unix nano second - should be unpredictable enough
	token := fmt.Sprintf("%v", time.Now().UnixNano())
	token = token[len(token)-6:]

	err = cred.SetPasswordResetToken(sr.DB, token)
	if err != nil {
		return &response, controller.HandlerError{500, "Failed to initiate password reset", err}
	}

    response.ResetPasswordToken = token
	return &response, nil
}
