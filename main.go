package main

import(
	"log"
	"net/http"
	"github.com/obedtandadjaja/auth-go/api"
)

func main() {
	http.Handle("/token", logRequestMiddleware(http.HandlerFunc(api.Token)))

	log.Println(http.ListenAndServe(":8000", nil))
}

func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s\n", r)
		next.ServeHTTP(w, r)
	})
}
