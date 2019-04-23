package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"bytes"

	"github.com/obedtandadjaja/auth-go/controller"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

type App struct {
	Router *mux.Router
}

func (app *App) Initialize(host string, port string, user string, password string, dbName string) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
		return
	}
	app.runMigration(db)

	app.Router = mux.NewRouter()
	sharedResources := &controller.SharedResources{ DB: db }
	app.initializeRoutes(sharedResources)
}

func (app *App) initializeRoutes(sr *controller.SharedResources) {
	for _, route := range routes {
		app.Router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(logRequestMiddleware(controller.Handler{sr, route.HandlerFunc}))
	}
}

func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// java StringBuilder equivalent
		var buffer bytes.Buffer

		buffer.WriteString(fmt.Sprintf("%v %v %v\n", r.Method, r.URL, r.Proto))
		buffer.WriteString(fmt.Sprintf("Host: %v\n", r.Host))

		// print header
		for name, headers := range r.Header {
			for _, header := range headers {
				buffer.WriteString(fmt.Sprintf("%v: %v | ", name, header))
			}
		}

		// if post then print form body
		if r.Method == "POST" {
			r.ParseForm()
			buffer.WriteString(r.Form.Encode())
		}

		log.Println(buffer.String())

		next.ServeHTTP(w, r)
	})
}

func (app *App) runMigration(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}

	migration, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
		return
	}

	migration.Steps(2)
}

func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}
