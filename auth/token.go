package auth

import (
	"log"
	"time"
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret-key")

var users = map[string]string {
	"user1@gmail.com": "password1",
	"user2@gmail.com": "password2",
}

type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claim struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type TokenResponse struct {
	Jwt string
}

func Token(w http.ResponseWriter, r *http.Request) {
	var cred Credential

	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := users[cred.Email]

	if !ok || expectedPassword != cred.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claim{
		Email: cred.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Fatal(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := TokenResponse{ Jwt: tokenString }
	json.NewEncoder(w).Encode(response)
}
