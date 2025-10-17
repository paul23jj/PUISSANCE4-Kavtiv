package pion

import (
	"fmt"
	"time"
)

const (
	Rows = 6
	Cols = 7
)

type Board struct {
	Grid [Rows][Cols]int
}

type Game struct {
	Board     Board  `json:"board"`
	Player    int    `json:"player"`
	LastState string `json:"state"`
}

func NewGame() *Game {
	return &Game{Player: 1}
}

func (g *Game) PlayMove(col int) error {
	ok, r, c := g.Board.AnimateDrop(col, g.Player) // 👈 ici on utilise la version animée
	if !ok {
		return fmt.Errorf("colonne %d pleine ou invalide", col)
	}

	g.LastState = g.Board.GameState(r, c, g.Player)

	if g.LastState == "En cours" {
		if g.Player == 1 {
			g.Player = 2
		} else {
			g.Player = 1
		}
	}

	return nil
}

func (g *Game) GetState() interface{} {
	return g
}

// 🌟 Version animée de Drop : on montre le pion tomber
func (b *Board) AnimateDrop(col, player int) (bool, int, int) {
	if col < 0 || col >= Cols {
		return false, -1, -1
	}

	for r := 0; r < Rows; r++ {
		// efface la position précédente
		if r > 0 {
			b.Grid[r-1][col] = 0
		}
		b.Grid[r][col] = player
		printBoard(b)
		time.Sleep(150 * time.Millisecond)

		// si on touche un pion en dessous ou le bas du plateau
		if r == Rows-1 || b.Grid[r+1][col] != 0 {
			return true, r, col
		}
	}
	return false, -1, -1
}

// 🌟 Fonction d’affichage du plateau dans le terminal
func printBoard(b *Board) {
	fmt.Print("\033[H\033[2J") // Efface la console (effet animation)
	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols; c++ {
			switch b.Grid[r][c] {
			case 0:
				fmt.Print(". ")
			case 1:
				fmt.Print("🔴 ")
			case 2:
				fmt.Print("🟡 ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

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

func (b *Board) IsFull() bool {
	for c := 0; c < Cols; c++ {
		if b.Grid[0][c] == 0 {
			return false
		}
	}
	return true
}

func (b *Board) GameState(lastRow, lastCol, player int) string {
	if lastRow >= 0 && lastCol >= 0 && b.IsWin(lastRow, lastCol) {
		return fmt.Sprintf("Victoire joueur %d", player)
	}
	if b.IsFull() {
		return "Match nul"
	}
	return "En cours"
}

func (b *Board) GridSlice() [][]int {
	g := make([][]int, Rows)
	for r := 0; r < Rows; r++ {
		g[r] = make([]int, Cols)
		for c := 0; c < Cols; c++ {
			g[r][c] = b.Grid[r][c]
		}
	}
	return g
}

func ExampleUsage() {
	game := NewGame()

	for {
		var col int
		fmt.Printf("Joueur %d, choisis une colonne (0-6): ", game.Player)
		fmt.Scan(&col)
		err := game.PlayMove(col)
		if err != nil {
			fmt.Println(err)
			continue
		}
		printBoard(&game.Board)
		fmt.Println(game.LastState)
		if game.LastState != "En cours" {
			break
		}
	}
}
