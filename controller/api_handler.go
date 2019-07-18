package controller

import (
	"database/sql"
	"log"
	"net/http"
)

type HttpError interface {
	error
	Status() int
}

type HandlerError struct {
	Code        int
	Err         error
	OriginalErr error
}

func (error HandlerError) Status() int {
	return error.Code
}

func (error HandlerError) Error() string {
	return error.Err.Error()
}

func (error HandlerError) OriginalError() (string, bool) {
	if error.OriginalErr == nil {
		return "", false
	}

	return error.OriginalErr.Error(), true
}

type SharedResources struct {
	DB  *sql.DB
	Env string
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
		case HandlerError:
			if originalErr, ok := e.OriginalError(); ok {
				log.Printf("ERROR: %s\n", originalErr)
			}

			log.Printf("HTTP %d - %s\n", e.Status(), e)

			// on prod error codes >= 500 should not be returned
			if h.SharedResources.Env == "production" && e.Status() >= 500 {
				http.Error(w, "Internal Server Error", e.Status())
			} else {
				http.Error(w, e.Error(), e.Status())
			}
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
