package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("."))
	log.Println("Serving at http://0.0.0.0:8000")
	log.Fatal(http.ListenAndServe(":8000", fs))
}
