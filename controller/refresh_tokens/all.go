package refresh_tokens

import "github.com/obedtandadjaja/auth-go/models/refresh_token"

type GetAllRequest struct {
	Uuid string `json:"credential_uuid"`
}

type GetAllResponse struct {
	RefreshTokens []*refresh_token.RefreshToken `json:"refresh_tokens"`
}
