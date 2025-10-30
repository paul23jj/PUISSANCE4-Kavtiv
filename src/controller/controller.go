package controller

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"puissance4/pion"
	"sync"
)

// --- 🌟 Variables globales ---
var ScoreJoueur1 int
var ScoreJoueur2 int
var gameInstance *pion.Game
var gameMu sync.Mutex

// GameSnapshot est une copie en lecture seule de l'état du jeu
type GameSnapshot struct {
	Grid   [6][7]int
	Player int
	State  string
}

// PlayMoveSafe joue un coup en protégeant l'accès concurrent à l'instance de jeu
func PlayMoveSafe(col int) error {
	gameMu.Lock()
	if gameInstance == nil {
		gameMu.Unlock()
		return fmt.Errorf("jeu non initialisé")
	}

	// Si la partie est déjà terminée (victoire / nul) et quelqu'un tente de jouer,
	// on réinitialise la partie AVANT d'appliquer le nouveau coup.
	if gameInstance.LastState != "En cours" {
		*gameInstance = *pion.NewGame()
	}

	// état avant le coup (devrait être "En cours" après l'éventuelle réinitialisation)
	prevState := gameInstance.LastState

	// joue le coup (met à jour LastState)
	err := gameInstance.PlayMove(col)
	state := gameInstance.LastState

	// n'incrémente le score que si on est passé de "En cours" → "Victoire ..."
	if prevState == "En cours" {
		if state == "Victoire joueur 1" {
			ScoreJoueur1++
		} else if state == "Victoire joueur 2" {
			ScoreJoueur2++
		}
	}

	gameMu.Unlock()

	// ne pas réinitialiser ici : la réinitialisation se fera lorsque quelqu'un essaiera
	// de jouer à nouveau après la victoire (logique gérée ci‑dessus).
	return err
}

// Snapshot retourne une copie sûre de l'état courant du jeu
func Snapshot() GameSnapshot {
	gameMu.Lock()
	defer gameMu.Unlock()
	snap := GameSnapshot{}
	if gameInstance == nil {
		return snap
	}
	// copie la grille
	for r := 0; r < 6; r++ {
		for c := 0; c < 7; c++ {
			snap.Grid[r][c] = int(gameInstance.Board.Grid[r][c])
		}
	}
	snap.Player = gameInstance.Player
	snap.State = gameInstance.LastState
	return snap
}

// ResetGame réinitialise la partie courante (thread-safe)
func ResetGame() {
	gameMu.Lock()
	defer gameMu.Unlock()
	if gameInstance != nil {
		*gameInstance = *pion.NewGame()
	}
}

// --- ⚙️ Fonctions utilitaires pour les templates ---
var funcMap = template.FuncMap{
	"inSlice": func(value string, list []string) bool {
		for _, item := range list {
			if item == value {
				return true
			}
		}
		return false
	},
	"seq": func(n int) []int {
		s := make([]int, n)
		for i := 0; i < n; i++ {
			s[i] = i
		}
		return s
	},
	"eq": func(a, b int) bool {
		return a == b
	},
}

// --- 🧱 Rendu d’un template avec recherche automatique ---
func renderTemplate(w http.ResponseWriter, filename string, data interface{}) {
	baseDir, _ := os.Getwd()

	// 🔍 Chemins possibles
	paths := []string{
		filepath.Join(baseDir, "src", "template", filename),
		filepath.Join(baseDir, "template", filename),
	}

	var tmplPath string
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			tmplPath = p
			break
		}
	}

	// ⚠️ Si introuvable
	if tmplPath == "" {
		http.Error(w, "Template introuvable : "+filename, http.StatusInternalServerError)
		fmt.Println("❌ Template non trouvé dans :", paths)
		return
	}

	fmt.Println("📂 Template utilisé :", tmplPath)

	tmpl := template.Must(
		template.New(filepath.Base(filename)).
			Funcs(funcMap).
			ParseFiles(tmplPath),
	)

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// --- 🧩 Injection du jeu ---
func SetGame(g *pion.Game) {
	gameInstance = g
}

// RenderGrid rend la template server-side de la grille en lisant les cookies
func RenderGrid(w http.ResponseWriter, r *http.Request) {
	snap := Snapshot()

	// defaults
	pawn1 := "/images/pawn1.svg"
	pawn2 := "/images/pawn2.svg"

	// lire cookies et appliquer si présents
	if c, err := r.Cookie("pionJoueur1"); err == nil && c.Value != "" {
		pawn1 = "/images/" + c.Value
	}
	if c, err := r.Cookie("pionJoueur2"); err == nil && c.Value != "" {
		pawn2 = "/images/" + c.Value
	}

	data := map[string]interface{}{
		"Grid":     snap.Grid,
		"Player":   snap.Player,
		"State":    snap.State,
		"PawnImg1": pawn1,
		"PawnImg2": pawn2,
		"Score1":   ScoreJoueur1,
		"Score2":   ScoreJoueur2,
	}

	renderTemplate(w, "grille.html", data)
}

// --- 🏠 Page d’accueil ---
func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Récupère les pions déjà pris
	taken := []string{}
	if c, err := r.Cookie("pionJoueur1"); err == nil && c.Value != "" {
		taken = append(taken, c.Value)
	}
	if c, err := r.Cookie("pionJoueur2"); err == nil && c.Value != "" {
		taken = append(taken, c.Value)
	}

	type ViewData struct {
		Title      string
		Message    string
		Grid       [6][7]int
		Player     int
		State      string
		Name1      string
		Name2      string
		PawnImg1   string
		PawnImg2   string
		TakenPawns []string
		Score1     int
		Score2     int
		BoardHTML  template.HTML
	}
	grid := make([][]int, 6)
	for i := range grid {
		grid[i] = make([]int, 7)
	}

	// utilise un snapshot sécurisé du jeu
	snap := Snapshot()

	// récupère images et noms depuis cookies si présents
	pawn1 := "/images/pawn1.svg"
	pawn2 := "/images/pawn2.svg"
	name1 := "Joueur 1"
	name2 := "Joueur 2"

	if c, err := r.Cookie("pionJoueur1"); err == nil && c.Value != "" {
		pawn1 = "/images/" + c.Value
	}
	if c, err := r.Cookie("pionJoueur2"); err == nil && c.Value != "" {
		pawn2 = "/images/" + c.Value
	}
	if c, err := r.Cookie("nomJoueur1"); err == nil && c.Value != "" {
		name1 = c.Value
	}
	if c, err := r.Cookie("nomJoueur2"); err == nil && c.Value != "" {
		name2 = c.Value
	}

	vd := ViewData{
		Title:      "Puissance 4",
		Message:    "Bienvenue sur la page d'accueil 🎉",
		Name1:      name1,
		Name2:      name2,
		Grid:       snap.Grid,
		Player:     snap.Player,
		State:      snap.State,
		PawnImg1:   pawn1,
		PawnImg2:   pawn2,
		TakenPawns: taken,
		Score1:     ScoreJoueur1,
		Score2:     ScoreJoueur2,
		BoardHTML:  buildBoardHTML(grid),
	}

	renderTemplate(w, "index.html", vd)
}

// --- 👥 Sélection du joueur ---
func Joueur(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Erreur formulaire: "+err.Error(), http.StatusBadRequest)
			return
		}

		joueur := r.FormValue("joueur")
		if joueur == "" {
			joueur = "1"
		}

		name := r.FormValue("name")
		pionChoisi := r.FormValue("pion")

		// --- Upload d’image ---
		file, _, err := r.FormFile("customImage")
		var imgName string
		if err == nil && file != nil {
			defer file.Close()

			imagesDir := filepath.Join("src", "images")
			os.MkdirAll(imagesDir, 0755)

			imgName = joueur
			outPath := filepath.Join(imagesDir, imgName)

			outFile, ferr := os.Create(outPath)
			if ferr == nil {
				defer outFile.Close()
				io.Copy(outFile, file)
			} else {
				imgName = pionChoisi
			}
		} else {
			imgName = fmt.Sprintf(pionChoisi)
		}

		// --- Cookies joueurs ---
		if joueur == "2" {
			http.SetCookie(w, &http.Cookie{Name: "nomJoueur2", Value: name, Path: "/"})
			http.SetCookie(w, &http.Cookie{Name: "pionJoueur2", Value: imgName, Path: "/"})
		} else {
			http.SetCookie(w, &http.Cookie{Name: "nomJoueur1", Value: name, Path: "/"})
			http.SetCookie(w, &http.Cookie{Name: "pionJoueur1", Value: imgName, Path: "/"})
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := map[string]string{
		"Title":   "Sélection du joueur",
		"Message": "Choisis ton pion et entre ton prénom",
	}
	renderTemplate(w, "player.html", data)
}

func buildBoardHTML(grid [][]int) template.HTML {
	html := "<table class='board'>"
	for _, row := range grid {
		html += "<tr>"
		for _, cell := range row {
			html += fmt.Sprintf("<td>%d</td>", cell)
		}
		html += "</tr>"
	}
	html += "</table>"
	return template.HTML(html)
}

// --- 🔄 Réinitialise la partie ---
func Reset(w http.ResponseWriter, r *http.Request) {
	if gameInstance != nil {
		*gameInstance = *pion.NewGame()
	}

	for _, c := range []string{"nomJoueur1", "nomJoueur2", "pionJoueur1", "pionJoueur2"} {
		http.SetCookie(w, &http.Cookie{Name: c, Value: "", Path: "/", MaxAge: -1})
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type Move struct {
	Col, Row, Player int
}

type ViewData struct {
	Players  map[int]string
	LastMove *Move
	Winner   int
}
// --- 🎵 Rappeurs autorisés = noms de fichiers mp3 dans /sounds (sans extension) ---
var allowedRappers = map[string]bool{
	"sdm": true, "kaaris": true, "booba": true, "naps": true, "damso": true,
	"tiako": true, "leto": true, "plk": true, "niska": true, "jul": true,
}

func stripExt(filename string) string {
	i := len(filename) - 1
	for i >= 0 && filename[i] != '.' && filename[i] != '/' && filename[i] != '\\' {
		i--
	}
	if i >= 0 && filename[i] == '.' {
		return filename[:i]
	}
	return filename
}

func imgToRapperID(img string) string {
	id := stripExt(img)       // "booba.png" -> "booba"
	if allowedRappers[id] {
		return id
	}
	return "booba" // fallback safe
}

// lastPlayerFromSnap : qui a joué le DERNIER coup ?
func lastPlayerFromSnap(player int, state string) int {
	// Dans ton moteur:
	// - si "En cours" => g.Player a déjà été togglé ⇒ le dernier joueur = l'autre
	// - si "Victoire ..." => g.Player N'A PAS été togglé ⇒ le dernier joueur = g.Player (le vainqueur)
	if state == "En cours" {
		if player == 1 { return 2 }
		return 1
	}
	// si Victoire ou Match nul:
	return player
}

// winnerFromState : 0 / 1 / 2 selon LastState
func winnerFromState(state string) int {
	switch state {
	case "Victoire joueur 1":
		return 1
	case "Victoire joueur 2":
		return 2
	default:
		return 0
	}
}
