package main

import (
	"fmt"
	"net/http"
	"power4/router"
)

func main() {
	// Charge le routeur
	r := router.New()

	fmt.Println("🚀 Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
