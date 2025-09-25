package routeur

import (
	"encoding/json"
	"net/http"

	"../pion"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

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

		err = game.PlayMove(data.Col)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(game.GetState())
	})

	// Route pour récupérer l'état du plateau
	mux.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(game.GetState())
	})

	return mux
}
