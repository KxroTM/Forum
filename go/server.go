package forum

import (
	"log"
	"net/http"
	"strings"
	"time"
)

var ResetPasswordMap = make(map[string]string)
var URL string

type DataStruct struct {
	User             User
	UserTarget       User
	RecommendedUser  RecommendedUser
	AllUsers         []User
	Post             Post
	AllPosts         []Post
	Comment          Comment
	AllComments      []Comment
	Notification     Notification
	AllNotifications []Notification
	Error            error
}

type RecommendedUser struct {
	RecommendedUsers []User
	Reason           []string
}

var AllData DataStruct

func CreateRoute(w http.ResponseWriter, r *http.Request, url string) {
	URL = url + "/"
	ResetPasswordMap[URL] = "valid"
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)

	AllData = GetAllDatas()

	// Si l'utilisateur est déjà connecté, on le redirige vers la page d'accueil
	data, _ := getSessionData(r)
	if data.User.Email != "" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		check := r.FormValue("save")

		connected, err := LoginUser(Db, email, password)

		if err == nil && connected {
			user := GetAccount(Db, email)

			if check == "" {
				err := createSessionCookie(w, SessionData{
					User: Session{
						UUID:      user.User_id,
						Email:     user.Email,
						Username:  user.Username,
						Role:      user.Role,
						ColorMode: "light",
					},
				}, 24*time.Hour)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			} else {
				err := createSessionCookie(w, SessionData{
					User: Session{
						UUID:      user.User_id,
						Email:     user.Email,
						Username:  user.Username,
						Role:      user.Role,
						ColorMode: "light",
					},
				}, 730*time.Hour)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}

			clientIP := r.RemoteAddr
			err = AccountLog(clientIP + "  ==>  " + email)
			if err != nil {
				log.Println(err)
			}
			http.Redirect(w, r, "/accueil", http.StatusSeeOther)
			return

		} else {
			AllData.Error = err
			err := LoginError.ExecuteTemplate(w, "loginerror.html", AllData)
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
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	err = AccountLog(clientIP + "  <==  " + UserSession.Email)
	if err != nil {
		log.Println(err)
	}
	LogoutUser()
	deleteSessionCookie(w)
	http.Redirect(w, r, "/accueil", http.StatusSeeOther)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)

	AllData = GetAllDatas()

	// Si l'utilisateur est déjà connecté, on le redirige vers la page d'accueil
	data, _ := getSessionData(r)
	if data.User.Email != "" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		passwordcheck := r.FormValue("passwordcheck")
		err := SignUpUser(Db, username, email, password, passwordcheck)

		if err == nil {
			http.Redirect(w, r, "/connexion", http.StatusSeeOther)
			return
		} else {
			AllData.Error = err
			err := RegisterError.ExecuteTemplate(w, "registererror.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	p := "Register page"
	err = Register.ExecuteTemplate(w, "register.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)

	AllData = GetAllDatas()

	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

	data, _ := getSessionData(r)
	if data.User.ColorMode == "dark" {
		err := DarkHome.ExecuteTemplate(w, "home.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = Home.ExecuteTemplate(w, "home.html", AllData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)
	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 3 || !strings.HasPrefix(parts[2], "@") {
		err := Error404.ExecuteTemplate(w, "error.html", "Invalid URL")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	username := strings.TrimPrefix(parts[2], "@")

	AllData := GetAllDatas()
	AllData.UserTarget = GetAccountByUsername(Db, username)

	if AllData.UserTarget == (User{}) {
		err := ErrorUser.ExecuteTemplate(w, "errorUser.html", "Invalid URL")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = Profile.ExecuteTemplate(w, "profile.html", AllData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusNotFound)
	p := "Page not found"
	err = Error404.ExecuteTemplate(w, "error.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
