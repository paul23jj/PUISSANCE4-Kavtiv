package routeur

import (
	"encoding/json"
	"net/http"

	"puissance4/pion"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	// 🌟 Création d'une instance du jeu (si nécessaire)
	game := pion.NewGame() // 🌟 nouvelle ligne — à adapter selon ton package "pion"

	// Serve files statiques (CSS/JS/images)
	// expose /static/ -> src/static/ and /images/ -> src/images/
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("src/images"))))

	// Route pour jouer un coup
	mux.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		var data struct {
			Col int `json:"col"`
		}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "JSON invalide", http.StatusBadRequest)
			return
		}

		// 🌟 Appel de la méthode PlayMove sur l'instance du jeu
		err = game.PlayMove(data.Col) // 🌟 nouvelle ligne
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 🌟 Encode l'état du jeu après le coup
		json.NewEncoder(w).Encode(game.GetState()) // 🌟 nouvelle ligne
	})

	// Route pour récupérer l'état du plateau
	mux.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		// 🌟 Encode l'état actuel du jeu
		json.NewEncoder(w).Encode(game.GetState()) // 🌟 nouvelle ligne
	})

	return mux
}
