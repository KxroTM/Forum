// Fichier main.go
package main

import (
	"database/sql"
	"log"
	"net/http"
	"os/exec"

	forum "forum/go"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./src/styles"))))
	http.HandleFunc("/", forum.NotFoundHandler)
	http.HandleFunc("/connexion", forum.LoginPage)
	http.HandleFunc("/inscription", forum.RegisterPage)
	http.HandleFunc("/deconnexion", forum.LogoutPage)
	http.HandleFunc("/google-login", forum.GoogleLoginPage)
	http.HandleFunc("/google-callback", forum.GoogleCallback)
	http.HandleFunc("/github-login", forum.GitHubLoginPage)
	http.HandleFunc("/github-callback", forum.GitHubCallback)
	http.HandleFunc("/accueil", forum.HomePage)
	http.HandleFunc("/profile/", forum.ProfilePage)
	http.HandleFunc("/creer-un-post", forum.CreatePostPage)

	err := openLink()
	if err != nil {
		log.Fatal("Erreur lors de l'ouverture du lien:", err)
		return
	}
	http.ListenAndServe(":8080", nil)
}

func init() {
	var err error
	forum.Db, err = sql.Open("sqlite3", "./database/db.sql")
	if err != nil {
		log.Fatal("Erreur lors de l'ouverture de la base de données:", err)
		return
	}
}

func openLink() error {
	cmd := exec.Command("cmd", "/c", "start", "http://localhost:8080/accueil")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
