package pion

import "fmt"

const (
	Rows = 6 // 6 lignes
	Cols = 7 // 7 colonnes
)

type Board struct {
	Grid [Rows][Cols]int // 0 vide, 1 joueur1, 2 joueur2
}

// Ajouter un pion dans une colonne
func (b *Board) Drop(col, player int) (bool, int, int) {
	if col < 0 || col >= Cols {
		return false, -1, -1
	}
	for r := Rows - 1; r >= 0; r-- {
		if b.Grid[r][col] == 0 {
			b.Grid[r][col] = player
			return true, r, col
		}
	}
	return false, -1, -1 // colonne pleine
}

// Vérifie si un joueur a gagné à partir d’une position
func (b *Board) IsWin(r, c int) bool {
	player := b.Grid[r][c]
	if player == 0 {
		return false
	}
	dirs := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for _, d := range dirs {
		count := 1
		count += b.checkLine(r, c, d[0], d[1], player)
		count += b.checkLine(r, c, -d[0], -d[1], player)
		if count >= 4 {
			return true
		}
	}
	return false
}

func (b *Board) checkLine(r, c, dr, dc, player int) int {
	count := 0
	r += dr
	c += dc
	for r >= 0 && r < Rows && c >= 0 && c < Cols && b.Grid[r][c] == player {
		count++
		r += dr
		c += dc
	}
	return count
}

// Vérifie si le plateau est plein
func (b *Board) IsFull() bool {
	for c := 0; c < Cols; c++ {
		if b.Grid[0][c] == 0 {
			return false
		}
	}
	return true
}

// Renvoie l’état du jeu
func (b *Board) GameState(lastRow, lastCol, player int) string {
	if lastRow >= 0 && lastCol >= 0 && b.IsWin(lastRow, lastCol) {
		return fmt.Sprintf("Victoire joueur %d", player)
	}
	if b.IsFull() {
		return "Match nul"
	}
	return "En cours"
}

// Exemple d’utilisation
func ExampleUsage() {
	var board Board

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
