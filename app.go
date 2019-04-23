package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

type App struct {
	Router *mux.Router
}

func (app *App) Initialize(username, password, dbName string) {
	connectionString := fmt.Sprintf("user=%s, password=%s dbname=%s", username, password, dbName)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	sharedResources := &handler.SharedResources{
		DB: db,
	}

	app.runMigration()

	app.Router = mux.NewRouter()
	app.initializeRoutes(sharedResources)
}

func (app *App) initializeRoutes(SharedResources *SharedResources) {
	for _, route := range routes {
		app.Router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(Handler{sharedResources, logRequestMiddleware(route.HandlerFunc)})
	}
}

func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s\n", r)
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
