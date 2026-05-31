package main

import (
	"log"
	"net/http"
)

func main() {
	mux := Rute()

	log.Println("Serverul Bubble Planet rulează pe http://localhost:9000 ...")

	err := http.ListenAndServe(":9000", mux)
	if err != nil {
		log.Fatalf("Eroare la pornirea serverului: %v", err)
	}
}
