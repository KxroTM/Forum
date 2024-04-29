// Fichier main.go
package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os/exec"

	forum "forum/go"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./src/styles"))))

	http.HandleFunc("/", forum.NotFoundHandler)
	http.HandleFunc("/google-login", forum.GoogleLoginPage)
	http.HandleFunc("/google-callback", forum.GoogleCallback)

	http.ListenAndServe(":8080", nil)
	openLink()
}

func init() {
	var err error
	forum.Db, err = sql.Open("sqlite3", "./database/db.db")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture de la base de données:", err)
		return
	}
	UpdateDb(forum.Db)
}

func UpdateDb(Db *sql.DB) {
	forum.UpdateUserDb(Db)
	forum.UpdatePostDb(Db)
	forum.UpdateCommentDb(Db)
}

func openLink() {
	cmd := exec.Command("cmd", "/c", "start", "http://localhost:8080/accueil")
	cmd.Run()
}
