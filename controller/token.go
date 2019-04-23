package controller

import (
	"net/http"
	"encoding/json"

	"github.com/obedtandadjaja/auth-go/models/credential"
	"github.com/obedtandadjaja/auth-go/auth/jwt"
)

type TokenRequest struct {
	Identifier    string `json:"identifier"`
	Password string `json:"password"`
}

type TokenResponse struct {
	jwt string
}

func Token(sr *SharedResources, w http.ResponseWriter, r *http.Request) error {
	var request TokenRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	credential, err := credential.FindBy(sr.DB, "identifier", request.Identifier)

	if credential.Password != request.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}

	tokenString, err := jwt.Generate(request.Identifier)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	response := TokenResponse{ jwt: tokenString }
	json.NewEncoder(w).Encode(response)

	return nil
}
