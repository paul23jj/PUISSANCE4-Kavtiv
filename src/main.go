package main

import (
	"fmt"
	"log"
	"net/http"
	"puissance4/routeur"
)

func main() {
	r := routeur.New()
	fmt.Println("🚀 Serveur démarré sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
