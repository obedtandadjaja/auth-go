package credential

import (
	"time"
	"github.com/Masterminds/squirrel"
)

type Credential struct {
	Id int32
	Identifier string
	Password string
	Subject string
	LastSignedIn time
	CreatedAt Time
	UpdatedAt Time
}

func All() []Credential {
	db, err := dbx.Open("mysql", "obedt:@localhost/auth")

	q := db.NewQuery("SELECT * FROM credentials LIMIT 10")

	var credentials []Credential
	err = q.All(&credentials)
}
