package main

import (
	"net/http"
)

func Rute() *http.ServeMux {
	m := http.NewServeMux()

	// 1. Ruta fixa pentru Service Worker
	m.HandleFunc("GET /sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "./webserver/pwa/sw.js")
	})

	// 2. Ruta abstracta dinamica: prinde si radacina "/" si orice string de tip "/ID_BULA"
	m.HandleFunc("/{rol...}", home)

	return m
}
