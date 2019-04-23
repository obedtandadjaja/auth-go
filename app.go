package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/obedtandadjaja/auth-go/controller"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialize(host string, port string, user string, password string, dbName string) {
	err := app.initializeDB(host, port, user, password, dbName)
	if err != nil {
		log.Fatal(err)
		return
	}
	app.runMigration()

	app.Router = mux.NewRouter()
	sharedResources := &controller.SharedResources{ DB: app.DB }
	app.initializeRoutes(sharedResources)
}

func (app *App) initializeDB(host string, port string, user string, password string, dbName string) error {
	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbName,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
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
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(requestDump))

		next.ServeHTTP(w, r)
	})
}

func (app *App) runMigration() {
	driver, err := postgres.WithInstance(app.DB, &postgres.Config{})
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
