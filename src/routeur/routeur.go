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

	// 🌟 Création d'une instance du jeu (si nécessaire)
	game := pion.NewGame() // 🌟 nouvelle ligne — à adapter selon ton package "pion"

	// Passe l'instance du jeu au controller pour rendu server-side
	controller.SetGame(game)

	// Serve files statiques (CSS/JS/images)
	// expose /static/ -> src/static/ and /images/ -> src/images/
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))
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
				http.Error(w, "Colonne invalide", http.StatusBadRequest)
				return
			}
			err = game.PlayMove(col)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// Rediriger vers la page d'accueil pour affichage HTML
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Sinon on essaie le JSON (API)
		var data struct {
			Col int `json:"col"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "JSON invalide", http.StatusBadRequest)
			return
		}
		err = game.PlayMove(data.Col)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
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
