package main

import (
	"fmt"
	"net/http"

	"./board"
	"./routeur"
)

func main() {
	// Charge le routeur
	r := routeur.New()

	fmt.Println("🚀 Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

func testBoard() {
	var board board.Board
	// Joueur 1 joue en colonne 3
	ok, r, c := board.Drop(3, 1)
	if ok {
		fmt.Println(board.GameState(r, c, 1)) // → "En cours"
	}
	// Joueur 2 joue en colonne 3
	ok, r, c = board.Drop(3, 2)
	if ok {
		fmt.Println(board.GameState(r, c, 2)) // → "En cours"
	}
}
