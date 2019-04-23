package credential

import (
	"time"
	"database/sql"
)

type Credential struct {
	Id           int32
	Identifier   string // can be email/username/phone
	Password     string
	Subject      string
	LastSignedIn time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IpAddress    string
}

func All(db *sql.DB) ([]*Credential, error) {
	rows, err := db.Query("select * from credentials")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	credentials := []*Credential{}
	for rows.Next() {
		credential, err := buildFromRows(rows)
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

func buildFromRow(row *sql.Row) (*Credential, error) {
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

func buildFromRows(rows *sql.Rows) (*Credential, error) {
	var credential Credential

	err := rows.Scan(
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
