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
	// Si le formulaire HTML envoie une colonne via POST, on redirige vers /play pour traitement
	if r.Method == http.MethodPost {
		r.ParseForm()
		// rediriger vers /play en POST standard (routeur gère form ou JSON)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Préparer les données pour le template : grille et état du jeu
	type ViewData struct {
		Title   string
		Message string
		Grid    [][]int
		Player  int
		State   string
	}

	vd := ViewData{Title: "Accueil", Message: "Bienvenue sur la page d'accueil 🎉", Grid: make([][]int, 6)}
	if gameInstance != nil {
		// copier la grille
		for r := 0; r < 6; r++ {
			row := make([]int, 7)
			for c := 0; c < 7; c++ {
				row[c] = int(gameInstance.Board.Grid[r][c])
			}
			vd.Grid[r] = row
		}
		vd.Player = gameInstance.Player
		vd.State = gameInstance.LastState
	}

	renderTemplate(w, "index.html", vd) // Affiche le template index.html avec les données
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

		name := r.FormValue("name")
		pion := r.FormValue("pion")

		// Si un fichier a été uploadé sous le champ 'photo', on l'enregistre dans src/images/pawn{N}.ext
		file, header, err := r.FormFile("photo")
		if err == nil && file != nil {
			defer file.Close()

			// S'assure que le dossier images existe
			imagesDir := "src/images"
			os.MkdirAll(imagesDir, 0755)

			// devine l'extension à partir du nom de fichier uploadé
			ext := filepath.Ext(header.Filename)
			if ext == "" {
				ext = ".png"
			}

			outPath := filepath.Join(imagesDir, fmt.Sprintf("pawn%s%s", pion, ext))

			outFile, err := os.Create(outPath)
			if err == nil {
				defer outFile.Close()
				io.Copy(outFile, file)
			}
		}

		// Enregistrer le choix dans des cookies simples (pour usage client)
		http.SetCookie(w, &http.Cookie{Name: "nomJoueur", Value: name, Path: "/"})
		http.SetCookie(w, &http.Cookie{Name: "pionJoueur", Value: pion, Path: "/"})

		data := map[string]string{
			"Title":   "Joueur enregistré",
			"Message": "Merci " + name + ". Tu as choisi le pion " + pion + ".",
		}
		renderTemplate(w, "player.html", data)
		return
	}

	data := map[string]string{
		"Title":   "Sélection du joueur",
		"Message": "Choisis ton pion et entre ton prénom",
	}
	renderTemplate(w, "player.html", data)
}
