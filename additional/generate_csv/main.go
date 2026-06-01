package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Valori implicite
	var lungimeString int64 = 17 // Lungimea din exemplul tău (JQGNJGTLLWZZGNMCM)
	numeFisier := "output.csv"
	var numarLinii int64 = 5000 // Numărul implicit de înregistrări

	arguments := os.Args

	// Citim argumentele din terminal sequentially
	if len(arguments) > 1 {
		lungimeString, _ = strconv.ParseInt(arguments[1], 10, 64)
	}
	if len(arguments) > 2 {
		numeFisier = arguments[2]
	}
	if len(arguments) > 3 {
		numarLinii, _ = strconv.ParseInt(arguments[3], 10, 64)
	}

	fmt.Printf("Pornire generare: %d linii, ID lungime %d -> %s\n", numarLinii, lungimeString, numeFisier)

	// 1. Creăm fișierul pe disc
	fisier, err := os.Create(numeFisier)
	if err != nil {
		fmt.Printf("Eroare la crearea fișierului: %v\n", err)
		return
	}
	defer fisier.Close()

	writer := csv.NewWriter(fisier)
	defer writer.Flush()

	// Scrie Antetul (Header)
	header := []string{"ID_Bula", "Data_Creare", "Status"}
	if err := writer.Write(header); err != nil {
		fmt.Printf("Eroare la scrierea antetului: %v\n", err)
		return
	}

	startChar := byte('A')
	timestampAcum := time.Now().Format("2006-01-02 15:04:05")

	// 2. Bucla principală care generează numărul de linii cerut
	var linie int64
	for linie = 0; linie < numarLinii; linie++ {

		// Generăm ID-ul unic pentru linia curentă
		idAleatoriu := ""
		var i int64 = 1
		for {
			myRand := random(0, 26)
			idAleatoriu += string(startChar + byte(myRand))
			if i == lungimeString {
				break
			}
			i++
		}

		// Scriem rândul direct în buffer-ul CSV
		dateRow := []string{idAleatoriu, timestampAcum, "activ"}
		if err := writer.Write(dateRow); err != nil {
			fmt.Printf("Eroare la scrierea liniei %d: %v\n", linie, err)
			return
		}
	}

	fmt.Printf("Succes! Am generat %d linii în [%s].\n", numarLinii, numeFisier)
}
