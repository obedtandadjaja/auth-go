package api

import (
	"net/http"
	"encoding/json"
	"github.com/obedtandadjaja/auth-go/auth/jwt"
)

var users = map[string]string {
	"user1@gmail.com": "password1",
	"user2@gmail.com": "password2",
}

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Jwt string
}

func Token(w http.ResponseWriter, r *http.Request) {
	var request TokenRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := users[request.Email]

	if !ok || expectedPassword != request.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString, err := jwt.Generate(request.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := TokenResponse{ Jwt: tokenString }
	json.NewEncoder(w).Encode(response)
}
