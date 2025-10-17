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
	soundsPath := filepath.Join(wd, "sounds") // ✅ <-- dossier sons ajouté ici

	fmt.Println("🔍 Test chemins :")
	fmt.Println("   Static =", staticPath)
	fmt.Println("   Images =", imagesPath)
	fmt.Println("   Sounds =", soundsPath) // ✅ log pour vérifier

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(imagesPath))))
	mux.Handle("/sounds/", http.StripPrefix("/sounds/", http.FileServer(http.Dir(soundsPath)))) // ✅ <-- nouvelle ligne

	fmt.Println("✅ Static mount: /static/ ->", staticPath)
	fmt.Println("✅ Static mount: /images/ ->", imagesPath)
	fmt.Println("✅ Static mount: /sounds/ ->", soundsPath) // ✅ confirmation sons

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
			// construire HTML minimal affichant la grille et boutons
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintln(w, "<html><head><title>Grille Puissance 4</title></head><body>")
			fmt.Fprintln(w, "<h1>Grille (jeu)</h1>")
			// utiliser Snapshot thread-safe du controller
			snap := controller.Snapshot()
			for i := 0; i < len(snap.Grid); i++ {
				fmt.Fprint(w, "<div style='display:flex'>")
				for j := 0; j < len(snap.Grid[i]); j++ {
					v := snap.Grid[i][j]
					cell := "&nbsp;"
					switch v {
					case 1:
						cell = "X"
					case 2:
						cell = "O"
					}
					fmt.Fprintf(w, "<div style='width:36px;height:36px;border:1px solid #333;display:flex;align-items:center;justify-content:center;margin:2px;'>%s</div>", cell)
				}
				fmt.Fprintln(w, "</div>")
			}

			// formulaire de jeu
			fmt.Fprintln(w, "<form method='post' action='/grille'>")
			for c := 0; c < 7; c++ {
				fmt.Fprintf(w, "<button type='submit' name='col' value='%d'>%d</button>", c, c+1)
			}
			fmt.Fprintln(w, "</form>")
			fmt.Fprintln(w, "<p><a href='/'>Retour accueil</a></p>")
			fmt.Fprintln(w, "</body></html>")
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
			snap := controller.Snapshot()
			switch snap.State {
			case "Victoire joueur 1":
				controller.ScoreJoueur1++
			case "Victoire joueur 2":
				controller.ScoreJoueur2++
			}
			http.Redirect(w, r, "/grille", http.StatusSeeOther)
		default:
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	})
	return mux
}
