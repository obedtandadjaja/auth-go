package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claim struct {
	Identifier string `json:"identifier"`
	jwt.StandardClaims
}

func Generate(identifier string) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claim{
		Identifier: identifier,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey())
	if err != nil {
		return "", fmt.Errorf("error exchanging jwt token")
	}

	return tokenString, nil
}

func Verify(tokenString string) error {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return secretKey(), nil
		})

	if token.Claims.ExpiresAt < time.Now().Unix() {
		return fmt.Errorf("token has expired")
	}

	return err
}

func secretKey() []byte {
	return []byte(os.Getenv("SECRET_KEY"))
}
