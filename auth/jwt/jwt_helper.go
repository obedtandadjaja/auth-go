package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type RefreshTokenClaim struct {
	CredentialUuid string `json:"credential_uuid"`
	SessionUuid    string `json:"session_uuid"`
	jwt.StandardClaims
}

type AccessTokenClaim struct {
	CredentialUuid string `json:"credential_uuid"`
	jwt.StandardClaims
}

func GenerateRefreshToken(credentialUuid, sessionUuid string) (string, error) {
	expirationTime := time.Now().Add(10 * 24 * time.Hour)
	claims := &RefreshTokenClaim{
		CredentialUuid: credentialUuid,
		SessionUuid:    sessionUuid,
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

func GenerateAccessToken(credentialUuid string) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &AccessTokenClaim{
		CredentialUuid: credentialUuid,
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

func VerifyRefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&RefreshTokenClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return secretKey(), nil
		})
	if err != nil {
		return "", err
	}

	return token.Claims.(*RefreshTokenClaim).SessionUuid, nil
}

func VerifyAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessTokenClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return secretKey(), nil
		})
	if err != nil {
		return "", err
	}

	return token.Claims.(*AccessTokenClaim).CredentialUuid, nil
}

func secretKey() []byte {
	return []byte(os.Getenv("SECRET_KEY"))
}
