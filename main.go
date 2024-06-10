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
	forum.InvalidAllMail(forum.Db)
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./src/styles"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	http.HandleFunc("/", forum.NotFoundHandler)
	http.HandleFunc("/connexion", forum.LoginPage)
	http.HandleFunc("/inscription", forum.RegisterPage)
	http.HandleFunc("/deconnexion", forum.LogoutPage)
	http.HandleFunc("/google-login", forum.GoogleLoginPage)
	http.HandleFunc("/google-callback", forum.GoogleCallback)
	http.HandleFunc("/github-login", forum.GitHubLoginPage)
	http.HandleFunc("/github-callback", forum.GitHubCallback)
	http.HandleFunc("/accueil", forum.HomePage)
	http.HandleFunc("/populaire", forum.PopulairePage)
	http.HandleFunc("/profile/", forum.ProfilePage)
	http.HandleFunc("/creer-un-post", forum.CreatePostPage)
	http.HandleFunc("/post/", forum.PostPage)
	http.HandleFunc("/posts", forum.PostsPage)
	http.HandleFunc("/changeColorMode", forum.ChangeColorMode)
	http.HandleFunc("/notifications", forum.NotificationsPage)
	http.HandleFunc("/forgot-password", forum.ForgotPasswordPage)
	http.HandleFunc("/likePost", forum.LikeLogique)
	http.HandleFunc("/dislikePost", forum.DislikeLogique)
	http.HandleFunc("/retweetPost", forum.RetweetLogique)
	http.HandleFunc("/followUser", forum.FollowLogique)
	http.HandleFunc("/categorie", forum.CategoriePage)
	http.HandleFunc("/report", forum.ReportLogique)
	http.HandleFunc("/forgot-password-success", forum.ForgotPasswordSuccessPage)
	http.HandleFunc("/reset-password", forum.ResetPasswordPage)
	http.HandleFunc("/reset-password-success", forum.PasswordResetSuccessPage)
	http.HandleFunc("/no-mail-found", forum.NoMailFoundPage)
	http.HandleFunc("/lien-expire", forum.ExpiredLinkPage)

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
