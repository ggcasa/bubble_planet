// go build -o ../bin/ferry .
// mkdir $HOME/testunit
// cd ~/testunit/
// ../proiecte/bubble_planet/additional/bin/generate_csv 17 cars.csv 1000000
// Pornire generare: 1000000 linii, ID lungime 17 -> cars.csv
// Succes! Am generat 1000000 linii în [cars.csv].
//
// ********
//
//	../proiecte/bubble_planet/additional/bin/csv2sqlite3 -csv=cars.csv -db=cars.db
//
// Pornire import: [cars.csv] -> [cars.db] (pure Go)...
// Am procesat deja 100000 linii...
// Am procesat deja 200000 linii...
// Am procesat deja 300000 linii...
// Am procesat deja 400000 linii...
// Am procesat deja 500000 linii...
// Am procesat deja 600000 linii...
// Am procesat deja 700000 linii...
// Am procesat deja 800000 linii...
// Am procesat deja 900000 linii...
// Am procesat deja 1000000 linii...
//
// Succes! Am importat 1000000 rânduri în [cars.db] în 29.904125209s.
//  *********
// sqlite3 cars.db "SELECT count(*) FROM masini;"
// 1000000
//
// ************
// ~/testunit$ ../proiecte/bubble_planet/additional/bin/ferry
// 2026/06/06 10:25:34 Serverul a pornit pe http://localhost:9977
//  ********
//
// curl -i  http://localhost:9977/?limit=100000 > cars.json
//   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
//  Dload  Upload   Total   Spent    Left  Speed
// 100 6393k    0 6393k    0     0  19.6M      0 --:--:-- --:--:-- --:--:-- 19.5M
//
//  *********
// ls -lh
// total 98M
// -rw-rw-r-- 1 ggcasa ggcasa  29M Jun  6 10:03 cars.csv
// -rw-r--r-- 1 ggcasa ggcasa  63M Jun  6 10:15 cars.db
// -rw-r--r-- 1 ggcasa ggcasa  32K Jun  6 10:25 cars.db-shm
// -rw-r--r-- 1 ggcasa ggcasa    0 Jun  6 10:15 cars.db-wal
// -rw-rw-r-- 1 ggcasa ggcasa 6.3M Jun  6 10:26 cars.json

// **********
//
// curl -i  http://localhost:9977/?limit=1000000 > million_bubbles.json
//
//	 % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
//	Dload  Upload   Total   Spent    Left  Speed
//
// 100 62.4M    0 62.4M    0     0  22.7M      0 --:--:--  0:00:02 --:--:-- 22.7M
// ggcasa@ggltp:~/testunit$ ls -lh million_bubbles.json
// -rw-rw-r-- 1 ggcasa ggcasa 63M Jun  6 10:42 million_bubbles.json
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
