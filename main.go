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
	err = UpdateDb(forum.Db)
	if err != nil {
		log.Fatal("Erreur lors de la mise à jour de la base de données:", err)
		return
	}
}

func UpdateDb(Db *sql.DB) error {
	err := forum.UpdateUserDb(Db)
	if err != nil {
		log.Println("Erreur lors de la mise à jour de la base de données des utilisateurs:", err)
		return err
	}
	err = forum.UpdatePostDb(Db)
	if err != nil {
		log.Println("Erreur lors de la mise à jour de la base de données des posts:", err)
		return err
	}
	err = forum.UpdateCommentDb(Db)
	if err != nil {
		log.Println("Erreur lors de la mise à jour de la base de données des commentaires:", err)
		return err
	}
	return nil
}

func openLink() error {
	cmd := exec.Command("cmd", "/c", "start", "http://localhost:8080/accueil")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
