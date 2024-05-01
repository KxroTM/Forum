package forum

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
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
	data, _ := getSessionData(r)
	if data.User.Role != hashPasswordSHA256("guest") {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}
	url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Fonction permettant d'avoir accès aux informations de l'utilisateur connecté via Google
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	IPsLog(clientIP + "  ==>  " + r.URL.Path)

	data, _ := getSessionData(r)
	if data.User.Role != hashPasswordSHA256("guest") {
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
		SignUpUser(Db, usertemp.Email, usertemp.Email, password)

		UserSession = GetAccount(usertemp.Email)
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

func generateStrongPassword() string {
	length := 12
	numBytes := length * 3 / 4
	randomBytes := make([]byte, numBytes)
	rand.Read(randomBytes)
	password := base64.URLEncoding.EncodeToString(randomBytes)
	password = password[:length]
	return password
}
