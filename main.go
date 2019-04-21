package main

import(
	"log"
	"net/http"
	"github.com/obedtandadjaja/auth-go/api"

	"database/sql"
    _ "github.com/lib/pq"
    "github.com/golang-migrate/migrate"
    "github.com/golang-migrate/migrate/database/postgres"
    _ "github.com/golang-migrate/migrate/source/file"
)

func main() {
	initDatabase()

	http.Handle("/token", logRequestMiddleware(http.HandlerFunc(api.Token)))

	log.Println(http.ListenAndServe(":8000", nil))
}

func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s\n", r)
		next.ServeHTTP(w, r)
	})
}

func initDatabase() {
	db, err := sql.Open("postgres", "postgres://obedt:@localhost/auth?sslmode=disable")
	if err != nil {
		log.Println(err)
		return
	}

    driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Println(err)
		return
	}

    m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		log.Println(err)
		return
	}

    m.Steps(2)
}
