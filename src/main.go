package main

import (
	"fmt"
	"net/http"

	"./routeur"
)

func main() {
	// Charge le routeur
	r := routeur.New()

	fmt.Println("🚀 Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
h