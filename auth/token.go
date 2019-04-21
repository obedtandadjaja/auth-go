package auth

import (
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret-key")

var users = map[string]string {
	"user1@gmail.com": "password1",
	"user2@gmail.com": "password2"
}

type Credential struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Claim struct {
	Username string `json:"email"`
	jwt.StandardClaims
}

func Token(w http.ResponseWriter, r *http.Request) {
	var cred Credential

	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := users[cred.Email]

	if !ok || expectedPassowrd != cred.Password {
		w.WriteHeader(http.StatusUnathorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claim{
		Username: cred.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Fatal(err)

		w.WriteHeader(http.StatusInternalServerEror)
		error
	}

	w.Write(tokenString)
}
