package forum

import (
	"net/http"
	"strings"
	"time"
)

var ResetPasswordMap = make(map[string]string)
var URL string

type DataStruct struct {
	User             User
	UserTarget       User
	AllUsers         []User
	Post             Post
	AllPosts         []Post
	Comment          Comment
	AllComments      []Comment
	Notification     Notification
	AllNotifications []Notification
}

var AllData DataStruct

func CreateRoute(w http.ResponseWriter, r *http.Request, url string) {
	URL = url + "/"
	ResetPasswordMap[URL] = "valid"
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	IPsLog(clientIP + "  ==>  " + r.URL.Path)
	updateUserSession(r)

	// Si l'utilisateur est déjà connecté, on le redirige vers la page d'accueil
	data, err := getSessionData(r)
	if err == nil || data.User.Email != "" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		check := r.FormValue("save")

		if LoginUser(Db, email, password) {
			user := GetAccount(email)

			if check == "" {
				createSessionCookie(w, SessionData{
					User: Session{
						UUID:     user.User_id,
						Email:    user.Email,
						Username: user.Username,
						Role:     user.Role,
					},
				}, 24*time.Hour)
			} else {
				createSessionCookie(w, SessionData{
					User: Session{
						UUID:     user.User_id,
						Email:    user.Email,
						Username: user.Username,
						Role:     user.Role,
					},
				}, 730*time.Hour)
			}

			clientIP := r.RemoteAddr
			AccountLog(clientIP + "  ==>  " + email)
			http.Redirect(w, r, "/accueil", http.StatusSeeOther)
			return

		} else {
			p := "Login page"
			err := LoginError.ExecuteTemplate(w, "loginerror.html", p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	p := "Login page"
	err = Login.ExecuteTemplate(w, "login.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LogoutPage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	IPsLog(clientIP + "  ==>  " + r.URL.Path)
	AccountLog(clientIP + "  <==  " + UserSession.Email)
	LogoutUser()
	deleteSessionCookie(w)
	http.Redirect(w, r, "/accueil", http.StatusSeeOther)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	IPsLog(clientIP + "  ==>  " + r.URL.Path)
	updateUserSession(r)

	AllData = GetAllDatas()

	err := Home.ExecuteTemplate(w, "home.html", AllData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	IPsLog(clientIP + "  ==>  " + r.URL.Path)
	updateUserSession(r)
	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 3 || !strings.HasPrefix(parts[2], "@") {
		err := Error.ExecuteTemplate(w, "error.html", "Invalid URL")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	username := strings.TrimPrefix(parts[2], "@")

	AllData := GetAllDatas()
	AllData.UserTarget = GetAccountByUsername(username)

	if AllData.UserTarget == (User{}) {
		err := ErrorUser.ExecuteTemplate(w, "errorUser.html", "Invalid URL")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err := Profile.ExecuteTemplate(w, "profile.html", AllData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	IPsLog(clientIP + "  ==>  " + r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	p := "Page not found"
	err := Error.ExecuteTemplate(w, "error.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
