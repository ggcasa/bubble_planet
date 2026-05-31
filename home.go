package main

import (
	"html/template"
	"net/http"
)

type Ticket struct {
	Name        string
	Description string
}

type Client struct {
	Name    string
	Rol     string
	Tickets []Ticket
}

var tp = template.Must(template.New("bule").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Bubble Planet</title>
</head>
<body>
    {{if .Tickets}}
        <!-- Ecran abstract de lucru (cazul true din switch) -->
        <h1>Home</h1>
        <h2>Name: {{.Name}} Rol: {{.Rol}}</h2>
        <ul>
            {{range .Tickets}}
            <li>Task: {{.Name}} Description: {{.Description}}</li>
            {{end}}
        </ul>
    {{else}}
        <!-- Formularul simplu afișat când URL-ul este gol "/" sau codul e invalid -->
        <input type="text" id="c">
        <button onclick="window.location.href='/'+document.getElementById('c').value">OK</button>
    {{end}}

    <script>
        if ('serviceWorker' in navigator) {
            navigator.serviceWorker.register('/sw.js');
        }
    </script>
</body>
</html>`))

// Functie simulata de verificare (pană se leaga cea reală din store/helpers)
func verificaCodInSistem(cod string) (Client, bool) {
	if cod == "123ABC" {
		return Client{
			Name: "66",
			Rol:  "Operator",
			Tickets: []Ticket{
				{Name: "Task1", Description: "Do that"},
				{Name: "Task2", Description: "Do this"},
			},
		}, true
	}
	return Client{}, false
}

func home(w http.ResponseWriter, r *http.Request) {
	// r.PathValue("rol") va fi:
	// - string gol "" dacă accesezi doar localhost:9000/
	// - stringul introdus (ex: "123ABC") daca accesezi localhost:9000/123ABC
	codIntroducs := r.PathValue("rol")

	client, codulEsteValid := verificaCodInSistem(codIntroducs)

	switch codulEsteValid {
	case true:

		tp.Execute(w, client)

	default:

		tp.Execute(w, Client{})
	}
}
