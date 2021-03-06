package session

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/obedtandadjaja/auth-go/models"
)

type Session struct {
	Id             int
	Uuid           string
	CredentialId   int
	IpAddress      sql.NullString
	UserAgent      sql.NullString
	LastAccessedAt time.Time
	CreatedAt      time.Time
	ExpiresAt      time.Time
}

func FindBy(db *sql.DB, fields map[string]interface{}) (*Session, error) {
	var findStatement []string
	var findValues []interface{}

	index := 0
	for k, v := range fields {
		index++
		findStatement = append(findStatement, fmt.Sprintf("%v = $%v", k, index))
		findValues = append(findValues, v)
	}

	sql := "select * from sessions where " + strings.Join(findStatement, " and ")

	return buildFromRow(db.QueryRow(sql, findValues...))
}

func (session *Session) Create(db *sql.DB) error {
	err := db.QueryRow(
		`insert into sessions
		 (credential_id, ip_address, user_agent, expires_at) values
		 ($1, $2, $3, $4) returning id, uuid`,
		session.CredentialId, session.IpAddress, session.UserAgent, session.ExpiresAt,
	).Scan(&session.Id, &session.Uuid)

	return err
}

func Where(db *sql.DB, fields map[string]interface{}) ([]*Session, error) {
    sessions := []*Session{}

	var whereStatement []string
	var whereValues []interface{}

	index := 0
	for k, v := range fields {
		index++
		whereStatement = append(whereStatement, fmt.Sprintf("%v = $%v", k, index))
		whereValues = append(whereValues, v)
	}

	sql := "select * from sessions where " + strings.Join(whereStatement, " and ")

    rows, err := db.Query(sql, whereValues...)
    if err != nil {
        return sessions, err
    }
    defer rows.Close()

    for rows.Next() {
        session, err := buildFromRow(rows)
        if err != nil {
            return sessions, err
        }
        sessions = append(sessions, session)
    }

    return sessions, nil
}

func (session *Session) Delete(db *sql.DB) error {
	_, err := db.Exec("delete from sessions where id = $1", session.Id)

	return err
}

func (session *Session) UpdateLastAccessedAt(db *sql.DB) error {
	_, err := db.Exec(`update sessions set last_accessed_at = now() where id = $1`, session.Id)

	return err
}

func buildFromRow(row models.ScannableObject) (*Session, error) {
	var session Session

	err := row.Scan(
		&session.Id,
		&session.Uuid,
		&session.CredentialId,
		&session.IpAddress,
		&session.UserAgent,
		&session.LastAccessedAt,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		fmt.Println(err)
		return &session, err
	}

	return &session, nil
}
