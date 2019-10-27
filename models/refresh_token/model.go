package refresh_token

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/obedtandadjaja/auth-go/models"
)

type RefreshToken struct {
	Id             int
	Uuid           string
	CredentialId   int
	IpAddress      string
	UserAgent      string
	LastAccessedAt time.Time
	CreatedAt      time.Time
	ExpiresAt      time.Time
}

func FindBy(db *sql.DB, fields map[string]interface{}) (*RefreshToken, error) {
	var findStatement []string
	var findValues []interface{}

	index := 0
	for k, v := range fields {
		index++
		findStatement = append(findStatement, fmt.Sprintf("%v = $%v", k, index))
		findValues = append(findValues, v)
	}

	sql := "select * from refresh_tokens where " + strings.Join(findStatement, " and ")

	return buildFromRow(db.QueryRow(sql, findValues...))
}

func (refreshToken *RefreshToken) Create(db *sql.DB) error {
	err := db.QueryRow(
		`insert into refresh_tokens
		 (credential_id, ip_address, user_agent, expires_at) values
		 ($1, $2, $3, $4) returning id, token`,
		refreshToken.CredentialId, refreshToken.IpAddress, refreshToken.UserAgent, refreshToken.ExpiresAt,
	).Scan(&refreshToken.Id, &refreshToken.Uuid)

	return err
}

func buildFromRow(row models.ScannableObject) (*RefreshToken, error) {
	var refreshToken RefreshToken

	err := row.Scan(
		&refreshToken.Id,
		&refreshToken.Uuid,
		&refreshToken.CredentialId,
		&refreshToken.IpAddress,
		&refreshToken.UserAgent,
		&refreshToken.LastAccessedAt,
		&refreshToken.CreatedAt,
		&refreshToken.ExpiresAt,
	)

	if err != nil {
		fmt.Println(err)
		return &refreshToken, err
	}

	return &refreshToken, nil
}
