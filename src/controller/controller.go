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

var ScoreJoueur1 int
var ScoreJoueur2 int

// renderTemplate est une fonction utilitaire pour afficher un template HTML avec des données dynamiques
// renderTemplate est une fonction utilitaire pour afficher un template HTML avec des données dynamiques
func renderTemplate(w http.ResponseWriter, filename string, data interface{}) {
	tmpl := template.Must(template.ParseFiles("template/" + filename)) // Charge le fichier template depuis le dossier "template"
	tmpl.Execute(w, data)                                              // Exécute le template et écrit le résultat dans la réponse HTTP
}

// instance du jeu (injectée depuis le routeur)
var gameInstance *pion.Game

// SetGame permet d'injecter une instance de jeu pour le rendu côté serveur
func SetGame(g *pion.Game) {
	gameInstance = g
}

// Home gère la page d'accueil
func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Récupère les pions déjà pris dans les cookies
	taken := []string{}
	if c, err := r.Cookie("pionJoueur1"); err == nil && c.Value != "" {
		taken = append(taken, c.Value)
	}
	if c, err := r.Cookie("pionJoueur2"); err == nil && c.Value != "" {
		taken = append(taken, c.Value)
	}

	// Préparer les données pour le template : grille et état du jeu
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
		Title:      "Accueil",
		Message:    "Bienvenue sur la page d'accueil 🎉",
		Name1: 	"Joueur 1",
		Name2: 	"Joueur 2",
		Grid:       make([][]int, 6),
		PawnImg1:   "/images/pawn1.svg",
		PawnImg2:   "/images/pawn2.svg",
		TakenPawns: taken,
		Score1:     11,
		Score2:     11,
	}
	// ...le reste du code...
	renderTemplate(w, "index.html", vd)
}

// Joueur affiche et gère le formulaire de sélection du joueur (prénom + pion)
func Joueur(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Accepte un formulaire multipart pour un upload optionnel de fichier
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, "Erreur formulaire: "+err.Error(), http.StatusBadRequest)
			return
		}

		joueur := r.FormValue("joueur")
		name := r.FormValue("name")
		pion := r.FormValue("pion")

		// Si un fichier a été uploadé sous le champ 'photo', on l'enregistre dans src/images/pawn{N}.ext
		file, header, err := r.FormFile("photo")
		var imgName string
		if err == nil && file != nil {
			defer file.Close()

			imagesDir := "src/images"
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
				// si erreur, retomber sur l'image choisie
				imgName = fmt.Sprintf("pawn%s.svg", pion)
			}
		} else {
			imgName = fmt.Sprintf("pawn%s.svg", pion)
		}

		// Enregistrer cookies spécifiques au joueur (1 ou 2)
		if joueur == "2" {
			http.SetCookie(w, &http.Cookie{Name: "nomJoueur2", Value: name, Path: "/"})
			// déterminer le nom d'image du pion choisi
			img := fmt.Sprintf("pawn%s.svg", pion)
			http.SetCookie(w, &http.Cookie{Name: "pionJoueur2", Value: img, Path: "/"})
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

// Reset réinitialise la partie et supprime les cookies joueurs
func Reset(w http.ResponseWriter, r *http.Request) {
	// Réinitialise l'instance du jeu
	if gameInstance != nil {
		*gameInstance = *pion.NewGame()
	}
	// Supprime les cookies joueurs
	for _, c := range []string{"nomJoueur1", "nomJoueur2", "pionJoueur1", "pionJoueur2"} {
		http.SetCookie(w, &http.Cookie{Name: c, Value: "", Path: "/", MaxAge: -1})
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
