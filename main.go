package main

import(
	"log"
	"net/http"
	"auth/token"
)

func main() {
	http.HandleFunc("/token", Token)
	// http.HandleFunc("/welcome", Welcome)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
