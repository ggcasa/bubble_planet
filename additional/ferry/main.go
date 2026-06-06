package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/glebarez/go-sqlite"
)

// Definim structura datelor care vor fi transformate în JSON
type Masina struct {
	ID_Bula string `json:"id_bula"`
	Marca   string `json:"marca"`
	Culoare string `json:"culoare"`
}

var db *sql.DB

func main() {
	var err error
	// Deschidem baza de date creată anterior
	db, err = sql.Open("sqlite", "cars.db")
	if err != nil {
		log.Fatalf("Eroare la deschiderea bazei de date: %v", err)
	}
	defer db.Close()

	// Activăm modul WAL pentru acces concurent sigur
	_, _ = db.Exec("PRAGMA journal_mode = WAL;")

	// Definim endpoint-ul API
	http.HandleFunc("/", masiniHandler)

	// Pornim serverul web pe portul 9977
	log.Println("Serverul a pornit pe http://localhost:9977")
	if err := http.ListenAndServe(":9977", nil); err != nil {
		log.Fatal(err)
	}
}

// func masiniHandler(w http.ResponseWriter, r *http.Request) {
// 	// Setăm header-ul ca browserul să știe că primește JSON
// 	w.Header().Set("Content-Type", "application/json")
// 	// Permitem accesul cross-origin (opțional, util în dezvoltare)
// 	w.Header().Set("Access-Control-Allow-Origin", "*")

// 	// Citim datele din SQLite (punem o limită de 100 pentru test)
// 	rows, err := db.Query("SELECT id_bula, marca, culoare FROM masini LIMIT 100")
// 	if err != nil {
// 		http.Error(w, "Eroare la citirea din baza de date", http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var listaMasini []Masina

// 	for rows.Next() {
// 		var m Masina
// 		if err := rows.Scan(&m.ID_Bula, &m.Marca, &m.Culoare); err != nil {
// 			continue
// 		}
// 		listaMasini = append(listaMasini, m)
// 	}

// 	// Transformăm slice-ul Go în JSON și îl trimitem pe rețea
// 	json.NewEncoder(w).Encode(listaMasini)
// }

// http://localhost:9977/?limit=5000
func masiniHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Preluăm parametrul "limit" din URL (ex: /?limit=250)
	limitaStr := r.URL.Query().Get("limit")

	// Valoarea implicită dacă utilizatorul nu specifică o limită
	limita := "100"
	if limitaStr != "" {
		limita = limitaStr
	}

	// Folosim o interogare securizată, transmițând limita ca parametru
	rows, err := db.Query("SELECT id_bula, marca, culoare FROM masini LIMIT ?", limita)
	if err != nil {
		http.Error(w, "Eroare la citirea din baza de date", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var listaMasini []Masina
	for rows.Next() {
		var m Masina
		if err := rows.Scan(&m.ID_Bula, &m.Marca, &m.Culoare); err != nil {
			continue
		}
		listaMasini = append(listaMasini, m)
	}

	json.NewEncoder(w).Encode(listaMasini)
}
