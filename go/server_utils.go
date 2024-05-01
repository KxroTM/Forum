package forum

import (
	"net/http"
	"time"
)

type Session struct {
	UUID     string
	Email    string
	Username string
	Role     string
}

type SessionData struct {
	User Session
}

func createSessionCookie(w http.ResponseWriter, data SessionData, hours time.Duration) {

	hashedUserRole := hashPasswordSHA256(data.User.Role)

	cookie := http.Cookie{
		Name:     "session",               // Nom du cookie
		Value:    hashedUserRole,          // Données du cookie
		HttpOnly: true,                    // Empêcher le JavaScript de lire le cookie
		Secure:   true,                    // Marquer le cookie comme sécurisé si vous utilisez HTTPS
		SameSite: http.SameSiteStrictMode, // Empêcher le navigateur d'envoyer le cookie avec les requêtes de site tiers
		Expires:  time.Now().Add(hours),   // Durée de validité du cookie, par exemple 24 heures
		Path:     "/",                     // Le cookie est valable pour tout le site
	}

	http.SetCookie(w, &cookie)
}

func getSessionData(r *http.Request) (SessionData, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return SessionData{
			User: Session{
				Role: hashPasswordSHA256("guest"),
			},
		}, err
	}

	data := SessionData{
		User: Session{
			Role: cookie.Value,
		},
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
		Expires:  time.Unix(0, 0),
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
}
