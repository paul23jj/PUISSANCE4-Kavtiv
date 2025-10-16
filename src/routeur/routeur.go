package routeur

import (
	"encoding/json"
	"fmt"
	"net/http"

	"puissance4/controller"
	"puissance4/pion"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/reset", controller.Reset)

	// 🌟 Création d'une instance du jeu (si nécessaire)
	game := pion.NewGame() // 🌟 nouvelle ligne — à adapter selon ton package "pion"

	// Passe l'instance du jeu au controller pour rendu server-side
	controller.SetGame(game)

	// Serve files statiques (CSS/JS/images)
	// expose /static/ -> src/static/ and /images/ -> src/images/
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))
	// expose /css/ -> src/css/ (le projet a un dossier src/css pour les styles)
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("src/css"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("src/images"))))

	// Route pour jouer un coup — accepte JSON (API) ou formulaire HTML (col)
	mux.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		// Priorité au formulaire HTML (r.FormValue)
		r.ParseForm()
		colStr := r.FormValue("col")
		if colStr != "" {
			// formulaire HTML : convertir en int
			var col int
			_, err := fmt.Sscanf(colStr, "%d", &col)
			if err != nil {
				// rediriger vers la page d'accueil avec message d'erreur
				http.Redirect(w, r, "/?err=Colonne+invalide", http.StatusSeeOther)
				return
			}
			err = game.PlayMove(col)
			if err != nil {
				// rediriger vers la page d'accueil avec message d'erreur
				http.Redirect(w, r, "/?err="+err.Error(), http.StatusSeeOther)
				return
			}
			// incrémenter score si victoire
			if game.LastState == "Victoire joueur 1" {
				controller.ScoreJoueur1++
			}
			if game.LastState == "Victoire joueur 2" {
				controller.ScoreJoueur2++
			}
			// Rediriger vers la page d'accueil pour affichage HTML
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Sinon on essaie le JSON (API)
		var data struct {
			Col int `json:"col"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "JSON invalide", http.StatusBadRequest)
			return
		}
		if err := game.PlayMove(data.Col); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Incrémente le score si victoire
		if game.LastState == "Victoire joueur 1" {
			controller.ScoreJoueur1++
		}
		if game.LastState == "Victoire joueur 2" {
			controller.ScoreJoueur2++
		}
		json.NewEncoder(w).Encode(game.GetState())
	})

	// Route pour récupérer l'état du plateau
	mux.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		// 🌟 Encode l'état actuel du jeu
		json.NewEncoder(w).Encode(game.GetState()) // 🌟 nouvelle ligne
	})

	// Pages templates basiques
	mux.HandleFunc("/", controller.Home)
	mux.HandleFunc("/joueur", controller.Joueur)

	return mux
}
