package board

import "fmt"

type Board struct {
	grid [6][7]int // 0 vide, 1 joueur1, 2 joueur2
}

func (b *Board) Drop(col int, player int) (bool, int, int) {
	return true, 0, col
}

func (b *Board) GameState(row int, col int, player int) string {
	return fmt.Sprintf("Player %d played at (%d, %d)", player, row, col)
}

func main() {

	var board Board

	// Joueur 1 joue
	ok, r, c := board.Drop(3, 1)
	if ok {
		fmt.Println(board.GameState(r, c, 1))
	}

	// Joueur 2 joue
	ok, r, c = board.Drop(3, 2)
	if ok {
		fmt.Println(board.GameState(r, c, 2))
	}
}
