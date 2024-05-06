package forum

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	ClientID:     "279057952766-sm8djuadlk75sh12oamo3fkb4q3r003l.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:8080/google-callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func GoogleLoginPage(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Fonction permettant d'avoir accès aux informations de l'utilisateur connecté via Google
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}

	data, _ := getSessionData(r)
	if data.User.Role != "guest" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}

	code := r.FormValue("code")

	token, err := googleOauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Échec de l'échange du code d'autorisation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfoResp, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, "Échec de la récupération des informations utilisateur depuis Google: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer userInfoResp.Body.Close()

	userInfoData, err := io.ReadAll(userInfoResp.Body)
	if err != nil {
		http.Error(w, "Échec de la lecture des données de réponse JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var usertemp User
	err = json.Unmarshal(userInfoData, &usertemp)
	if err != nil {
		http.Error(w, "Échec de l'analyse des données JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if FindAccount(usertemp.Email) {
		UserSession = GetAccount(usertemp.Email)
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
			http.Error(w, "Erreur lors de la création du cookie de session: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		password := generateStrongPassword() + "@1L"

		username := strings.Split(usertemp.Email, "@")

		err := SignUpUser(Db, username[0], usertemp.Email, password, password)

		if err != nil {
			http.Error(w, "Erreur lors de l'inscription de l'utilisateur: "+err.Error(), http.StatusInternalServerError)
			return
		}

		UserSession = GetAccount(usertemp.Email)
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
			http.Error(w, "Erreur lors de la création du cookie de session: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = AccountLog(clientIP + "  ==>  " + UserSession.Email)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/accueil", http.StatusSeeOther)
}
