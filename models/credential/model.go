package credential

import (
	"database/sql"
	"time"
	"strings"
	"fmt"

	"github.com/obedtandadjaja/auth-go/models"
	"github.com/obedtandadjaja/auth-go/auth/hash"

	"github.com/lib/pq"
)

type Credential struct {
	Id                 int
	Identifier         string // can be email/username/phone
	Password           sql.NullString
	Subject            sql.NullString
	LastSignedIn       pq.NullTime
	CreatedAt          pq.NullTime
	UpdatedAt          pq.NullTime
	IpAddress          sql.NullString
	FailedAttempts     int
	LockedUntil        pq.NullTime
	PasswordResetToken sql.NullString
}

func All(db *sql.DB) ([]*Credential, error) {
	rows, err := db.Query("select * from credentials")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	credentials := []*Credential{}
	for rows.Next() {
		credential, err := buildFromRow(rows)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, credential)
	}

	return credentials, nil
}

func FindBy(db *sql.DB, fields map[string]interface{}) (*Credential, error) {
	var findStatement []string
	for k, v := range fields {
		findStatement = append(findStatement, fmt.Sprintf("%v = %v", k, v))
	}

	return buildFromRow(db.QueryRow("select * from credentials where $1 limit 1", strings.Join(findStatement, " and ")))
}

func (credential *Credential) Create(db *sql.DB) error {
	hashValue, err := hash.HashPassword(credential.Password.String)
	if err != nil {
		return err
	}

	err = db.QueryRow(
		`insert into credentials
         (id, identifier, password, subject, last_signed_in, created_at, updated_at, ip_address) values
         (default, $1, $2, $3, $4, $5, $6, $7)
         returning id`,
		credential.Identifier, hashValue, credential.Subject, nil,
		time.Now(), time.Now(), credential.IpAddress,
	).Scan(&credential.Id)

	return err
}

func (credential *Credential) Update(db *sql.DB, fields map[string]interface{}) error {
	var updateStatement []string
	for k, v := range fields {
		updateStatement = append(updateStatement, fmt.Sprintf("%v = %v", k, v))
	}

	_, err := db.Exec("update credentials set $1 where id = $2", strings.Join(updateStatement, ","), credential.Id)

	return err
}

func (credential *Credential) IncrementFailedAttempt(db *sql.DB) error {
	_, err := db.Exec(`update credentials set failed_attempts = failed_attempts + 1
                       where id = $1 and failed_attempts = $2`, credential.Id, credential.FailedAttempts)

	return err
}

func (credential *Credential) Delete(db *sql.DB) error {
	_, err := db.Exec("delete from credentials where id = $1", credential.Id)

	return err
}

func (credential *Credential) UpdatePassword(db *sql.DB) error {
	hashValue, err := hash.HashPassword(credential.Password.String)
	if err != nil {
		return nil
	}

	_, err = db.Exec("update credentials set password = $1, password_reset_token = null where id = $2", hashValue, credential.Id)

	return err
}

func (credential *Credential) SetPasswordResetToken(db *sql.DB) error {
	hashValue, err := hash.HashPassword(fmt.Sprintf("%v", time.Now().Unix()))
	if err != nil {
		return nil
	}

	_, err = db.Exec("update credentials set password_reset_token = $1, where id = $2", hashValue, credential.Id)

	return err
}

func buildFromRow(row models.ScannableObject) (*Credential, error) {
	var credential Credential

	err := row.Scan(
		&credential.Id,
		&credential.Identifier,
		&credential.Password,
		&credential.Subject,
		&credential.LastSignedIn,
		&credential.CreatedAt,
		&credential.UpdatedAt,
		&credential.IpAddress,
		&credential.FailedAttempts,
		&credential.LockedUntil,
		&credential.PasswordResetToken,
	)

	if err != nil {
		return &credential, err
	}

	return &credential, nil
}
