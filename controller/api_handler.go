package controller

import (
	"net/http"
	"database/sql"
	"log"
)

type HttpError interface {
	error
	Status() int
}

type HandlerError struct {
	Code int
	Err  error
}

func (error HandlerError) Error() string {
	return error.Err.Error()
}

func (error HandlerError) Status() int {
	return error.Code
}

type SharedResources struct {
	DB *sql.DB
}

type Handler struct {
	SharedResources *SharedResources
	Handler         func(sr *SharedResources, w http.ResponseWriter, r *http.Request) error
}

// to satisfy http.Handler
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.Handler(h.SharedResources, w, r)

	if err != nil {
		switch e := err.(type) {
		case HttpError:
			log.Printf("HTTP %d - %s\n", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
