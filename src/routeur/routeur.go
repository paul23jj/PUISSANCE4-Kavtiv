package routeur

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"puissance4/controller"
	"puissance4/pion"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	// --- Affiche le répertoire de travail ---
	wd, _ := os.Getwd()
	fmt.Println("📂 Dossier de travail :", wd)

	// --- Sert les fichiers statiques ---
	staticPath := filepath.Join(wd, "static")
	imagesPath := filepath.Join(wd, "images")

	fmt.Println("🔍 Test chemins :")
	fmt.Println("   Static =", staticPath)
	fmt.Println("   Images =", imagesPath)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(imagesPath))))
	fmt.Println("✅ Static mount: /static/ ->", staticPath)
	fmt.Println("✅ Static mount: /images/ ->", imagesPath)

	// --- Routes principales ---
	mux.HandleFunc("/", controller.Home)
	mux.HandleFunc("/joueur", controller.Joueur)
	mux.HandleFunc("/reset", controller.Reset)

	// --- Jeu ---
	game := pion.NewGame()
	controller.SetGame(game)

	// --- API /play ---
	mux.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}
		var data struct {
			Col int `json:"col"`
		}
		json.NewDecoder(r.Body).Decode(&data)
		if err := game.PlayMove(data.Col); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if game.LastState == "Victoire joueur 1" {
			controller.ScoreJoueur1++
		} else if game.LastState == "Victoire joueur 2" {
			controller.ScoreJoueur2++
		}
		json.NewEncoder(w).Encode(game.GetState())
	})

	// --- Debug endpoint ---
	mux.HandleFunc("/__debug", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Working dir: %s\n", wd)
		files, _ := os.ReadDir(wd)
		fmt.Fprintf(w, "\nContenu du dossier :\n")
		for _, f := range files {
			fmt.Fprintf(w, " - %s\n", f.Name())
		}
	})

	return mux
}
