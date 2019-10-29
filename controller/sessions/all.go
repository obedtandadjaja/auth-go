package refresh_tokens

import "github.com/obedtandadjaja/auth-go/models/session"

type GetAllRequest struct {
	Uuid string `json:"credential_uuid"`
}

type GetAllResponse struct {
	RefreshTokens []*session.Session `json:"sessions"`
}
