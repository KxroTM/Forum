package forum

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

type Session struct {
	UUID      string
	Email     string
	Username  string
	Role      string
	ColorMode string
}

type SessionData struct {
	User Session
}

func createSessionCookie(w http.ResponseWriter, data SessionData, hours time.Duration) error {
	encodedData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	encodedString := base64.StdEncoding.EncodeToString(encodedData)

	cookie := http.Cookie{
		Name:     "session",               // Nom du cookie
		Value:    encodedString,           // Données du cookie (JSON encodé en base64)
		HttpOnly: true,                    // Empêcher le JavaScript de lire le cookie
		Secure:   true,                    // Marquer le cookie comme sécurisé si vous utilisez HTTPS
		SameSite: http.SameSiteStrictMode, // Empêcher le navigateur d'envoyer le cookie avec les requêtes de site tiers
		Expires:  time.Now().Add(hours),   // Durée de validité du cookie, par exemple 24 heures
		Path:     "/",                     // Le cookie est valable pour tout le site
	}

	http.SetCookie(w, &cookie)
	return nil
}

func getSessionData(r *http.Request) (SessionData, error) {
	var data SessionData

	cookie, err := r.Cookie("session")
	if err != nil {
		if AllData.ColorMode == "light" {
			return SessionData{
				User: Session{
					Role:      "guest",
					ColorMode: "light",
				},
			}, err
		} else {
			return SessionData{
				User: Session{
					Role:      "guest",
					ColorMode: "dark",
				},
			}, err
		}
	}

	decodedData, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		if AllData.ColorMode == "light" {
			return SessionData{
				User: Session{
					Role:      "guest",
					ColorMode: "light",
				},
			}, err
		} else {
			return SessionData{
				User: Session{
					Role:      "guest",
					ColorMode: "dark",
				},
			}, err
		}
	}

	err = json.Unmarshal(decodedData, &data)
	if err != nil {
		if AllData.ColorMode == "light" {
			return SessionData{
				User: Session{
					Role:      "guest",
					ColorMode: "light",
				},
			}, err
		} else {
			return SessionData{
				User: Session{
					Role:      "guest",
					ColorMode: "dark",
				},
			}, err
		}
	}
	if AllData.ColorMode == "light" {
		data.User.ColorMode = "light"
	} else {
		data.User.ColorMode = "dark"
	}
	return data, nil
}

func deleteSessionCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
}

func updateUserSession(r *http.Request) {
	data, _ := getSessionData(r)
	if data.User.Email == "" {
		return
	}
	if UserSession.Email == "" {
		UserSession = GetAccountById(Db, data.User.UUID)
	}
}
