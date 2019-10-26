package refresh_token

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/obedtandadjaja/auth-go/auth/secure_random"
	"github.com/obedtandadjaja/auth-go/models"
)

type RefreshToken struct {
	Id           int
	Token        string
	CredentialId int
	ExpiresAt    pq.NullTime
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
	token, err := secure_random.GenerateRandomString(500)

	err = db.QueryRow(
		`insert into refresh_tokens
		 (token, credential_id, expires_at) values
		 ($1, $2, $3) returning id, token`,
		token, refreshToken.CredentialId, refreshToken.ExpiresAt,
	).Scan(&refreshToken.Id, &refreshToken.Token)

	return err
}

func buildFromRow(row models.ScannableObject) (*RefreshToken, error) {
	var refreshToken RefreshToken

	err := row.Scan(
		&refreshToken.Id,
		&refreshToken.Token,
		&refreshToken.CredentialId,
		&refreshToken.ExpiresAt,
	)

	if err != nil {
		fmt.Println(err)
		return &refreshToken, err
	}

	return &refreshToken, nil
}
