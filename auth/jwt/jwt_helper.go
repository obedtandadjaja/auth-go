package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claim struct {
	CredentialId string `json:"credential_id"`
	jwt.StandardClaims
}

func Generate(credentialId string) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claim{
		CredentialId: credentialId,
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

func Verify(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return secretKey(), nil
		})

	if err != nil {
		return "", err
	}

	// this is already done via .Valid(); retaining this here for future examples
	if token.Claims.(*Claim).ExpiresAt < time.Now().Unix() {
		return "", fmt.Errorf("token has expired")
	}

	return token.Claims.(*Claim).CredentialId, nil
}

func secretKey() []byte {
	return []byte(os.Getenv("SECRET_KEY"))
}
