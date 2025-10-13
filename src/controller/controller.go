package controller

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// renderTemplate est une fonction utilitaire pour afficher un template HTML avec des données dynamiques
func renderTemplate(w http.ResponseWriter, filename string, data map[string]string) {
	tmpl := template.Must(template.ParseFiles("template/" + filename)) // Charge le fichier template depuis le dossier "template"
	tmpl.Execute(w, data)                                              // Exécute le template et écrit le résultat dans la réponse HTTP
}

// Home gère la page d'accueil
func Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title":   "Accueil",                           // Titre de la page
		"Message": "Bienvenue sur la page d'accueil 🎉", // Message affiché dans le template
	}
	renderTemplate(w, "index.html", data) // Affiche le template index.html avec les données
}

// About gère la page "À propos"
func About(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title":   "À propos",
		"Message": "Ceci est la page À propos ✨",
	}
	renderTemplate(w, "about.html", data) // Affiche le template about.html avec les données
}

// Contact gère la page de contact
func Contact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost { // Si le formulaire est soumis en POST
		// Récupération des données du formulaire
		name := r.FormValue("name") // Récupère le champ "name"
		msg := r.FormValue("msg")   // Récupère le champ "msg"

		data := map[string]string{
			"Title":   "Contact",
			"Message": "Merci " + name + " pour ton message : " + msg, // Message personnalisé après soumission
		}
		renderTemplate(w, "contact.html", data)
		return // On termine ici pour ne pas exécuter la partie GET
	}

	// Si ce n'est pas un POST, on affiche simplement le formulaire
	data := map[string]string{
		"Title":   "Contact",
		"Message": "Envoie-nous un message 📩",
	}
	renderTemplate(w, "contact.html", data)
}

// Player affiche et gère le formulaire de sélection de joueur (prénom + pion)
func Player(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Accepte un formulaire multipart pour un upload optionnel de fichier
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, "Erreur formulaire: "+err.Error(), http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		pawn := r.FormValue("pawn")

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

			outPath := filepath.Join(imagesDir, fmt.Sprintf("pawn%s%s", pawn, ext))

			outFile, err := os.Create(outPath)
			if err == nil {
				defer outFile.Close()
				io.Copy(outFile, file)
			}
		}

		// Enregistrer le choix dans des cookies simples (pour usage client)
		http.SetCookie(w, &http.Cookie{Name: "playerName", Value: name, Path: "/"})
		http.SetCookie(w, &http.Cookie{Name: "playerPawn", Value: pawn, Path: "/"})

		data := map[string]string{
			"Title":   "Joueur enregistré",
			"Message": "Merci " + name + ". Tu as choisi le pion " + pawn + ".",
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
