package credential

import (
	"database/sql"
	"time"

	"github.com/obedtandadjaja/auth-go/models"
	"github.com/obedtandadjaja/auth-go/auth/hash"

	"github.com/lib/pq"
)

type Credential struct {
	Id           int64
	Identifier   string // can be email/username/phone
	Password     sql.NullString
	Subject      sql.NullString
	LastSignedIn pq.NullTime
	CreatedAt    pq.NullTime
	UpdatedAt    pq.NullTime
	IpAddress    sql.NullString
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

func FindBy(db *sql.DB, fieldName string, arg interface{}) (*Credential, error) {
	return buildFromRow(db.QueryRow("select * from credentials where $1 = $2", fieldName, arg))
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
	)

	if err != nil {
		return &credential, err
	}

	return &credential, nil
}
