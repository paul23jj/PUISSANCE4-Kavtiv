package controller

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"puissance4/pion"
)

// --- 🌟 Variables globales ---
var ScoreJoueur1 int
var ScoreJoueur2 int
var gameInstance *pion.Game

// --- ⚙️ Fonctions utilitaires disponibles dans les templates ---
var funcMap = template.FuncMap{
	"inSlice": func(value string, list []string) bool {
		for _, item := range list {
			if item == value {
				return true
			}
		}
		return false
	},
}

// --- 🧱 Fonction pour charger et exécuter un template (chemin dynamique et sûr) ---
func renderTemplate(w http.ResponseWriter, filename string, data interface{}) {
	baseDir, _ := os.Getwd() // récupère le dossier courant (celui d’où tu lances "go run main.go")
	tmplPath := filepath.Join(baseDir, "src", "template", filename)

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

// --- 🧩 Injection de l’instance du jeu ---
func SetGame(g *pion.Game) {
	gameInstance = g
}

// --- 🏠 Page d’accueil ---
func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Récupère les pions déjà pris depuis les cookies
	taken := []string{}
	if c, err := r.Cookie("pionJoueur1"); err == nil && c.Value != "" {
		taken = append(taken, c.Value)
	}
	if c, err := r.Cookie("pionJoueur2"); err == nil && c.Value != "" {
		taken = append(taken, c.Value)
	}

	// Données envoyées au template
	type ViewData struct {
		Title      string
		Message    string
		Grid       [][]int
		Player     int
		State      string
		Name1      string
		Name2      string
		PawnImg1   string
		PawnImg2   string
		TakenPawns []string
		Score1     int
		Score2     int
	}

	vd := ViewData{
		Title:      "Puissance 4",
		Message:    "Bienvenue sur la page d'accueil 🎉",
		Name1:      "Joueur 1",
		Name2:      "Joueur 2",
		Grid:       make([][]int, 6),
		PawnImg1:   "/images/pawn1.svg",
		PawnImg2:   "/images/pawn2.svg",
		TakenPawns: taken,
		Score1:     ScoreJoueur1,
		Score2:     ScoreJoueur2,
	}

	renderTemplate(w, "index.html", vd)
}

// --- 👥 Sélection du joueur ---
func Joueur(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			http.Error(w, "Erreur formulaire: "+err.Error(), http.StatusBadRequest)
			return
		}

		joueur := r.FormValue("joueur")
		if joueur == "" {
			joueur = r.URL.Query().Get("joueur")
			if joueur == "" {
				joueur = "1"
			}
		}

		name := r.FormValue("name")
		pionChoisi := r.FormValue("pion")

		// --- Gestion upload d'image ---
		file, header, err := r.FormFile("photo")
		var imgName string
		if err == nil && file != nil {
			defer file.Close()

			imagesDir := filepath.Join("src", "images")
			os.MkdirAll(imagesDir, 0755)

			ext := filepath.Ext(header.Filename)
			if ext == "" {
				ext = ".png"
			}

			imgName = fmt.Sprintf("pawn%s%s", joueur, ext)
			outPath := filepath.Join(imagesDir, imgName)

			outFile, ferr := os.Create(outPath)
			if ferr == nil {
				defer outFile.Close()
				io.Copy(outFile, file)
			} else {
				imgName = fmt.Sprintf("pawn%s.svg", pionChoisi)
			}
		} else {
			imgName = fmt.Sprintf("pawn%s.svg", pionChoisi)
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

// --- 🔄 Réinitialise la partie ---
func Reset(w http.ResponseWriter, r *http.Request) {
	if gameInstance != nil {
		*gameInstance = *pion.NewGame()
	}

	// Supprime les cookies
	for _, c := range []string{"nomJoueur1", "nomJoueur2", "pionJoueur1", "pionJoueur2"} {
		http.SetCookie(w, &http.Cookie{Name: c, Value: "", Path: "/", MaxAge: -1})
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
