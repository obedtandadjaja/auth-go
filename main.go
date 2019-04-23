package main

import "os"

func main() {
	app := App{}
	app.Initialize(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	app.Run(":8080")
}
