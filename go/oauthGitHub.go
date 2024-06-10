package forum

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

func GitHubLoginPage(w http.ResponseWriter, r *http.Request) {

	var (
		oauthConf = &oauth2.Config{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  "http://localhost:8080/github-callback",
			Scopes:       []string{"user:email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
		}
		oauthStateString = "84983c60f7daadc1cb8698621f802c0d9f9a3c3c295c810748fb048115c186ec"
	)

	data, _ := getSessionData(r)
	if data.User.Role != "guest" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}
	url := oauthConf.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Fonction permettant d'avoir accès aux informations de l'utilisateur connecté via GitHub
func GitHubCallback(w http.ResponseWriter, r *http.Request) {

	var (
		oauthConf = &oauth2.Config{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  "http://localhost:8080/github-callback",
			Scopes:       []string{"user:email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
		}
		oauthStateString = "84983c60f7daadc1cb8698621f802c0d9f9a3c3c295c810748fb048115c186ec"
	)

	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}

	state := r.FormValue("state")
	if state != oauthStateString {
		log.Printf("L'état OAuth n'est pas valide: %s\n", state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Échec de l'échange du code d'autorisation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := oauthConf.Client(context.Background(), token)
	response, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		http.Error(w, "Erreur lors de l'appel de l'API GitHub: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var userEmails []struct {
		Email string `json:"email"`
	}
	err = json.NewDecoder(response.Body).Decode(&userEmails)
	if err != nil {
		http.Error(w, "Échec de l'analyse des données JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(userEmails) == 0 {
		http.Error(w, "Aucun email n'a été retourné par l'API GitHub", http.StatusInternalServerError)
		return
	}

	userEmail := userEmails[0].Email

	if FindAccount(Db, userEmail) {
		UserSession = GetAccount(Db, userEmail)
		err := createSessionCookie(w, SessionData{
			User: Session{
				UUID:      UserSession.User_id,
				Email:     UserSession.Email,
				Username:  UserSession.Username,
				Role:      UserSession.Role,
				ColorMode: "light",
			},
		}, 24*time.Hour)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		password := generateStrongPassword() + "@1L"

		username := strings.Split(userEmail, "@")

		err := SignUpUser(Db, username[0], userEmail, password, password)

		if err != nil {
			http.Error(w, "Erreur lors de l'inscription de l'utilisateur: "+err.Error(), http.StatusInternalServerError)
			return
		}

		UserSession = GetAccount(Db, userEmail)
		err = createSessionCookie(w, SessionData{
			User: Session{
				UUID:      UserSession.User_id,
				Email:     UserSession.Email,
				Username:  UserSession.Username,
				Role:      UserSession.Role,
				ColorMode: "light",
			},
		}, 24*time.Hour)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	err = AccountLog(clientIP + "  ==>  " + UserSession.Email)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/accueil", http.StatusSeeOther)
}
