// go build -o csv2sqlite3 main.go
// ./csv2sqlite3 -csv="marci.csv" -db="baza_mea.db"
package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	_ "github.com/glebarez/go-sqlite" // Driverul pure-Go
)

func main() {
	// Definirea flag-urilor
	// flag.String primește: numele flag-ului, valoarea implicită și o scurtă descriere
	csvFlag := flag.String("csv", "marci.csv", "Calea către fișierul CSV de intrare")
	dbFlag := flag.String("db", "baza_date.db", "Calea către fișierul bazei de date SQLite")

	// Parsează argumentele din linia de comandă
	flag.Parse()

	// Preluăm valorile din pointeri
	numeFisierCSV := *csvFlag
	numeBazaDate := *dbFlag

	start := time.Now()

	// 1. Deschidem fișierul CSV pentru citire
	csvFile, err := os.Open(numeFisierCSV)
	if err != nil {
		log.Fatalf("Eroare la deschiderea fișierului CSV [%s]: %v", numeFisierCSV, err)
	}
	defer csvFile.Close()

	// 2. Conectarea la SQLite
	db, err := sql.Open("sqlite", numeBazaDate)
	if err != nil {
		log.Fatalf("Eroare la conectarea la baza de date [%s]: %v", numeBazaDate, err)
	}
	defer db.Close()

	// Optimizări de viteză pentru SQLite
	_, _ = db.Exec("PRAGMA synchronous = OFF;")
	_, _ = db.Exec("PRAGMA journal_mode = MEMORY;")

	// 3. Crearea tabelului în baza de date
	createTableSQL := `CREATE TABLE IF NOT EXISTS masini (
		id_bula TEXT PRIMARY KEY,
		marca TEXT,
		culoare TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Eroare la crearea tabelului: %v", err)
	}

	reader := csv.NewReader(csvFile)

	// Sărim peste primul rând (antetul: ID_Bula, Marca, Culoare)
	_, err = reader.Read()
	if err != nil {
		log.Fatalf("Eroare la citirea antetului CSV: %v", err)
	}

	// 4. Începem o TRANZACȚIE pentru viteză masivă
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Eroare la inițierea tranzacției: %v", err)
	}

	// Pregătim query-ul de insert în avans
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO masini (id_bula, marca, culoare) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatalf("Eroare la pregătirea query-ului: %v", err)
	}
	defer stmt.Close()

	fmt.Printf("Pornire import: [%s] -> [%s] (pure Go)...\n", numeFisierCSV, numeBazaDate)

	var contor int64 = 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // Am ajuns la sfârșitul fișierului
		}
		if err != nil {
			log.Printf("Eroare la citirea liniei %d: %v", contor+2, err)
			continue
		}

		// record[0] = ID_Bula, record[1] = Marca, record[2] = Culoare
		_, err = stmt.Exec(record[0], record[1], record[2])
		if err != nil {
			log.Printf("Eroare la inserarea liniei %d: %v", contor+2, err)
			continue
		}

		contor++
		if contor%100000 == 0 {
			fmt.Printf("Am procesat deja %d linii...\n", contor)
		}
	}

	// 5. Salvarea tuturor datelor (Commit)
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Eroare la salvarea (commit-ul) datelor în SQLite: %v", err)
	}

	durata := time.Since(start)
	fmt.Printf("\nSucces! Am importat %d rânduri în [%s] în %v.\n", contor, numeBazaDate, durata)
}
