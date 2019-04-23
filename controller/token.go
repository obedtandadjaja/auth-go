package api

import (
	"net/http"
	"encoding/json"

	"github.com/obedtandadjaja/auth-go/models"
	"github.com/obedtandadjaja/auth-go/auth/jwt"

	"github.com/gorilla/mux"
)

type TokenRequest struct {
	Identifier    string `json:"identifier"`
	Password string `json:"password"`
}

type TokenResponse struct {
	jwt string
}

func Token(w http.ResponseWriter, r *http.Request) {
	var request TokenRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	credential, err := models.credential.FindBy(app.DB, "identifier", request.Identifier)

	if credential.Password != request.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString, err := jwt.Generate(request.Identifier)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := TokenResponse{ jwt: tokenString }
	json.NewEncoder(w).Encode(response)
}
