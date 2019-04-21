package jwt

import (
	"fmt"
	"time"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret-key")

type Claim struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func Generate(email string) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("error exchanging jwt token")
	}

	return tokenString, nil
}
