package controller

import (
	"html/controller"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, filename string, data map[string]string) {
	tmpl := template.Must(template.ParseFiles("template/" + filename))
	tmpl.Execute(w, data)
}

func Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title": "Acceuil",
		"Body":  "Bienvenue sur la page d'accueil!",
	}
	renderTemplate(w, "index.html", data)
}


func About(w, http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title": "À propos",
		"Body":  "Ceci est la page À propos ",
	}
	renderTemplate(w, "about.html", data)
}


func Contact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		name := r.FormValue("name")
		email := r.FormValue("msg")

		data := map[string]string{
			"title":   "Contact",
			"message": "Merci " + name + "pour ton message:" + msg,
		}
		renderTemplate(w, "contact.html", data)
		return
	}


	data := map[string]string{
		"Title": "Contact",
		"Body":  "Envoie nous ton message ",
	}
	renderTemplate(w, "contact.html", data)
}