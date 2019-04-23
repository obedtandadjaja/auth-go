package credential

import (
	"database/sql"

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
		credential, err := buildFromRows(rows)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, credential)
	}

	return credentials, nil
}

func FindBy(db *sql.DB, fieldName string, arg interface{}) (*Credential, error) {
	return buildFromRow(db.QueryRow("select * from credentials limit 1;"))
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
