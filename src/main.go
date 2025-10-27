package main

import (
	"fmt"
	"net/http"
	"puissance4/routeur"
)

func main() {
	r := routeur.New()
	fmt.Println("🚀 Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
