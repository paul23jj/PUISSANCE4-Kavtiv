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


    wd, _ := os.Getwd()
    fmt.Println("📂 Dossier de travail :", wd)

    staticPath := filepath.Join(wd, "static")
    imagesPath := filepath.Join(wd, "images")
    soundsPath := filepath.Join(wd, "soundtrack") 

    mux.Handle("/static/",     http.StripPrefix("/static/",     http.FileServer(http.Dir(staticPath))))
    mux.Handle("/images/",     http.StripPrefix("/images/",     http.FileServer(http.Dir(imagesPath))))
    mux.Handle("/soundtrack/", http.StripPrefix("/soundtrack/", http.FileServer(http.Dir(soundsPath))))

    mux.HandleFunc("/",        controller.Home)
    mux.HandleFunc("/joueur",  controller.Joueur)
    mux.HandleFunc("/reset",   controller.Reset)

    game := pion.NewGame()
    controller.SetGame(game)

    mux.HandleFunc("/grille", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
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
            if err := controller.PlayMoveSafe(col); err != nil {
                http.Redirect(w, r, "/grille", http.StatusSeeOther)
                return
            }
            http.Redirect(w, r, "/grille", http.StatusSeeOther)
        default:
            http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        }
    })

    return mux
}
