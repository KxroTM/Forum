package forum

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     "d36714255c2c1ce0ca0e",
		ClientSecret: "5a877183c40f8ef24e1356c159c116fc522c04de",
		RedirectURL:  "http://localhost:8080/github-callback",
		Scopes:       []string{"user:email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
	oauthStateString = "84983c60f7daadc1cb8698621f802c0d9f9a3c3c295c810748fb048115c186ec"
)

func GitHubLoginPage(w http.ResponseWriter, r *http.Request) {
	url := oauthConf.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Fonction permettant d'avoir accès aux informations de l'utilisateur connecté via GitHub
func GitHubCallback(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	IPsLog(clientIP + "  ==>  " + r.URL.Path)

	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("L'état OAuth n'est pas valide: %s\n", state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, _ := getSessionData(r)
	if data.User.Role != hashPasswordSHA256("guest") {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
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

	// Décodage des données JSON
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

	fmt.Println(userEmails)

	userEmail := userEmails[0].Email

	if FindAccount(userEmail) {
		UserSession = GetAccount(userEmail)
		createSessionCookie(w, SessionData{
			User: Session{
				UUID:     UserSession.User_id,
				Email:    UserSession.Email,
				Username: UserSession.Username,
				Role:     UserSession.Role,
			},
		}, 24*time.Hour)
	} else {

		password := hashPasswordSHA256(generateStrongPassword())
		SignUpUser(Db, userEmail, userEmail, password)

		UserSession = GetAccount(userEmail)
		createSessionCookie(w, SessionData{
			User: Session{
				UUID:     UserSession.User_id,
				Email:    UserSession.Email,
				Username: UserSession.Username,
				Role:     UserSession.Role,
			},
		}, 24*time.Hour)
	}

	AccountLog(clientIP + "  ==>  " + UserSession.Email)

	http.Redirect(w, r, "/accueil", http.StatusSeeOther)
}