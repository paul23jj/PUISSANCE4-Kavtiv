package routeur

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"puissance4/controller"
	"puissance4/pion"
	"strconv"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	// (la synchronisation du jeu est gérée dans controller)

	// --- Affiche le répertoire de travail ---
	wd, _ := os.Getwd()
	fmt.Println("📂 Dossier de travail :", wd)

	// --- Sert les fichiers statiques ---
	staticPath := filepath.Join(wd, "static")
	imagesPath := filepath.Join(wd, "images")
	soundsPath := filepath.Join(wd, "sounds") 
	fmt.Println("🔍 Test chemins :")
	fmt.Println("   Static =", staticPath)
	fmt.Println("   Images =", imagesPath)
	fmt.Println("   Sounds =", soundsPath) 

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(imagesPath))))
	mux.Handle("/sounds/", http.StripPrefix("/sounds/", http.FileServer(http.Dir(soundsPath))))

	fmt.Println("✅ Static mount: /static/ ->", staticPath)
	fmt.Println("✅ Static mount: /images/ ->", imagesPath)
	fmt.Println("✅ Static mount: /sounds/ ->", soundsPath)

	// --- Routes principales ---
	mux.HandleFunc("/", controller.Home)
	mux.HandleFunc("/joueur", controller.Joueur)
	mux.HandleFunc("/reset", controller.Reset)

	// --- Jeu ---
	game := pion.NewGame()
	controller.SetGame(game)

	// Handler simple qui affiche une grille HTML et permet de jouer via formulaire POST
	mux.HandleFunc("/grille", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// render via template server-side
			controller.RenderGrid(w, r)
		case http.MethodPost:
			if err := r.ParseForm(); err != nil {
				http.Error(w, "form error", http.StatusBadRequest)
				return
			}
			colStr := r.FormValue("col")
			col, err := strconv.Atoi(colStr)
			if err != nil {
				http.Error(w, "colonne invalide", http.StatusBadRequest)
				return
			}
			// thread-safe via controller.PlayMoveSafe
			if err := controller.PlayMoveSafe(col); err != nil {
				// rediriger vers GET avec message d'erreur simple
				http.Redirect(w, r, "/grille", http.StatusSeeOther)
				return
			}
			// mettre à jour les scores depuis le snapshot
			http.Redirect(w, r, "/grille", http.StatusSeeOther)
		default:
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	})
	return mux
}
