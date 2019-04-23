package main

import "net/http"

type Error interface {
	error
	Status() int
}

type HttpError struct {
	Code int
	Error error
}

func (httpError HttpError) Error() string {
	return httpError.Error.Error()
}

func (httpError HttpError) Status() int {
	return httpError.Code
}

type SharedResources struct {
	DB *sql.DB
}

type Handler struct {
	*SharedResources
	Handler func(sr *SharedResources, w http.ResponseWriter, r *http.Request) error
}

func (h Handler) serveHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.Handler(h.SharedResources, w, r)

	if err != nil {
		switch e :=err.(type) {
		case Error:
			log.Printf("HTTP %d - %s\n", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
