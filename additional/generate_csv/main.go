// ./additional/bin/generate_csv 17 document.csv 100000
// 17: Lungimea ID-ului generat aleatoriu
// "document.csv": Numele fișierului
// 10000: Numărul de linii
// ************************
// go run . 17 docume.csv 30000000
// Pornire generare: 30000000 linii, ID lungime 17 -> docume.csv
// Succes! Am generat 30000000 linii în [docume.csv].
// ************************
// ls -slh docume.csv
// 843M -rw-rw-r-- 1 ggcasa ggcasa 843M Jun  5 19:47 docume.csv
// ************************
// head -5 docume.csv
// ID_Bula,Marca,Culoare
// CIBEXZYVRSEIOPZCQ,DAC,Galben
// IHJZBIFCOZKMWALWM,Dacia,Verde
// ADMRWPRSYEPQAFHRV,Dacia,Negru
// HLZTNSAYGVXJDLWBC,Roman,Alb
package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {

	// Valori implicite
	var lungimeString int64 = 17
	numeFisier := "output.csv"
	var numarLinii int64 = 5000

	arguments := os.Args

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

	fisier, err := os.Create(numeFisier)
	if err != nil {
		fmt.Printf("Eroare la crearea fișierului: %v\n", err)
		return
	}
	defer fisier.Close()

	writer := csv.NewWriter(fisier)
	defer writer.Flush()

	// Am actualizat antetul coloanelor pentru a reflecta noile date
	header := []string{"ID_Bula", "Marca", "Culoare"}
	if err := writer.Write(header); err != nil {
		fmt.Printf("Eroare la scrierea antetului: %v\n", err)
		return
	}

	// Definirea opțiunilor pentru mărci și culori
	marci := []string{"Dacia", "Trabant", "Aro", "DAC", "Roman"}
	culori := []string{"Alb", "Negru", "Rosu", "Albastru", "Verde", "Galben", "Gri"}

	startChar := byte('A')

	var linie int64
	for linie = 0; linie < numarLinii; linie++ {
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

		// Selectare aleatorie din liste folosind lungimea feliei ca limită maximă
		marcaAleatorie := marci[random(0, len(marci))]
		culoareAleatorie := culori[random(0, len(culori))]

		// Scriere date în rând
		dateRow := []string{idAleatoriu, marcaAleatorie, culoareAleatorie}
		if err := writer.Write(dateRow); err != nil {
			fmt.Printf("Eroare la scrierea liniei %d: %v\n", linie, err)
			return
		}
	}

	fmt.Printf("Succes! Am generat %d linii în [%s].\n", numarLinii, numeFisier)
}
