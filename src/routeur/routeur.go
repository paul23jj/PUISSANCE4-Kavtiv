package routeur

import (
	"net/http"

	"../controller"
	"../pion"
)

// New crée et retourne un nouvel objet ServeMux configuré avec les routes de l'application
func New() *http.ServeMux {

	mux := http.NewServeMux() // Création d'un nouveau ServeMux, qui est un routeur simple pour les requêtes HTTP

	// On associe les chemins URL à des fonctions spécifiques du controller
	mux.HandleFunc("/", controller.Home)           // "/" correspond à la page d'accueil. Appelle la fonction Home du controller
	mux.HandleFunc("/about", controller.About)     // "/about" correspond à la page "À propos". Appelle la fonction About
	mux.HandleFunc("/contact", controller.Contact) // "/contact" correspond à la page de contact. Appelle la fonction Contact

	return mux // On retourne le routeur configuré
}
