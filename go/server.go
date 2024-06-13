package forum

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	AllCategories    []Category
	Categorie        Category
	Error            error
	ColorMode        string
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

	AllData = GetAllDatas(r)

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
						ColorMode: data.User.ColorMode,
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
						ColorMode: data.User.ColorMode,
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
			if AllData.ColorMode == "light" {
				AllData.Error = err
				err := LoginError.ExecuteTemplate(w, "loginerror.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			} else {
				AllData.Error = err
				err := DarkLoginError.ExecuteTemplate(w, "loginerror.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	}

	if AllData.ColorMode == "light" {
		p := "Login page"
		err = Login.ExecuteTemplate(w, "login.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		p := "Login page"
		err = DarkLogin.ExecuteTemplate(w, "login.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	AllData = GetAllDatas(r)

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
			if AllData.ColorMode == "light" {
				AllData.Error = err
				err := RegisterError.ExecuteTemplate(w, "registererror.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			} else {
				AllData.Error = err
				err := DarkRegisterError.ExecuteTemplate(w, "registererror.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	}

	if AllData.ColorMode == "light" {
		p := "Register page"
		err = Register.ExecuteTemplate(w, "register.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		p := "Register page"
		err = DarkRegister.ExecuteTemplate(w, "register.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)

	query := r.URL.RawQuery
	AllData = GetAllDatas(r)

	recherche := r.FormValue("recherche")
	if recherche != "" {
		AllData.AllPosts = GetPostBySearch(Db, recherche, GetAllPosts(Db))
	} else {
		AllData.AllPosts = GetAllPosts(Db)
	}

	if query == "pourtoi" {
		if AllData.User.Email == "" {
			http.Redirect(w, r, "/connexion", http.StatusSeeOther)
			return
		}
		AllData.AllPosts = ForYouPageAlgorithm(Db, UserSession.User_id)
	} else if query == "suivies" {
		if AllData.User.Email == "" {
			http.Redirect(w, r, "/connexion", http.StatusSeeOther)
			return
		}
		AllData.AllPosts = GetPostByFollowing(Db, UserSession.User_id)
	}

	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

	if AllData.User.Email == "" {
		if AllData.ColorMode == "light" {
			err = Home.ExecuteTemplate(w, "accueil.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := DarkHome.ExecuteTemplate(w, "accueil.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		AllData.AllNotifications = GetNotifications(Db, UserSession.User_id)

		if AllData.ColorMode == "light" {
			err = HomeLogged.ExecuteTemplate(w, "accueilLogged.html", AllData)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := DarkHomeLogged.ExecuteTemplate(w, "accueilLogged.html", AllData)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
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
		if AllData.ColorMode == "light" {
			err := Error404.ExecuteTemplate(w, "error.html", "Invalid URL")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		} else {
			err := DarkError.ExecuteTemplate(w, "error.html", "Invalid URL")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	username := strings.TrimPrefix(parts[2], "@")

	AllData = GetAllDatas(r)

	recherche := r.FormValue("recherche")
	if recherche != "" {
		AllData.AllPosts = GetPostBySearch(Db, recherche, AllData.AllPosts)
	} else {
		AllData.AllPosts, _ = GetAllPostsByUser(Db, AllData.UserTarget.User_id)
	}

	AllData.UserTarget = GetAccountByUsername(Db, username)
	AllData.AllPosts, _ = GetAllPostsByUser(Db, AllData.UserTarget.User_id)
	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

	query := r.URL.RawQuery

	if query == "posts" {
		AllData.AllPosts, _ = GetAllPostsByUser(Db, AllData.UserTarget.User_id)
	} else if query == "likes" {
		AllData.AllPosts = GetAllPostsByLike(Db, AllData.UserTarget.Username)
	} else if query == "retweets" {
		AllData.AllPosts = GetAllPostsByRetweet(Db, AllData.UserTarget.Username)
	} else if query == "comments" {
		// AllData.AllPosts = GetAllPostsByComment(Db, AllData.UserTarget.User_id)
	}

	if strings.Contains(AllData.UserTarget.FollowerList, UserSession.Username) {
		AllData.UserTarget.HeFollowed = true
	}

	if AllData.UserTarget == (User{}) {
		if AllData.ColorMode == "light" {
			err := ErrorUser.ExecuteTemplate(w, "errorUser.html", "Invalid URL")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		} else {
			err := DarkErrorUser.ExecuteTemplate(w, "errorUser.html", "Invalid URL")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	if AllData.User.Email == "" {
		if AllData.ColorMode == "light" {
			err = Profile.ExecuteTemplate(w, "profile.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err = DarkProfile.ExecuteTemplate(w, "profile.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		AllData.AllNotifications = GetNotifications(Db, UserSession.User_id)
		if AllData.ColorMode == "light" {
			err = ProfileLogged.ExecuteTemplate(w, "profileLogged.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err = DarkProfileLogged.ExecuteTemplate(w, "profileLogged.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

// Recherche a faire
func PostPage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)
	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 3 || !strings.HasPrefix(parts[2], "id=") {
		if AllData.ColorMode == "light" {
			err := Error404.ExecuteTemplate(w, "error.html", "Invalid URL")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		} else {
			err := DarkError.ExecuteTemplate(w, "error.html", "Invalid URL")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	id := strings.TrimPrefix(parts[2], "id=")

	AllData = GetAllDatas(r)
	AllData.Post, err = GetPost(Db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if AllData.Post == (Post{}) {
		if AllData.ColorMode == "light" {
			err := ErrorPost.ExecuteTemplate(w, "errorPost.html", "Invalid URL")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		} else {
			err := DarkErrorPost.ExecuteTemplate(w, "errorPost.html", "Invalid URL")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	if AllData.User.Email == "" {
		if AllData.ColorMode == "light" {
			err = LightPost.ExecuteTemplate(w, "post.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err = DarkPost.ExecuteTemplate(w, "post.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		AllData.AllNotifications = GetNotifications(Db, UserSession.User_id)
		if AllData.ColorMode == "light" {
			err = LightPostLogged.ExecuteTemplate(w, "postLogged.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err = DarkPostLogged.ExecuteTemplate(w, "postLogged.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

// Recherche a faire
func CreatePostPage(w http.ResponseWriter, r *http.Request) {
	data, _ := getSessionData(r)
	if data.User.Email == "" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}

	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)

	AllData = GetAllDatas(r)
	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

	if r.Method == http.MethodPost {
		postType := r.FormValue("post_type")
		if postType == "text" {
			title := r.FormValue("title")
			content := r.FormValue("text")
			categ := r.Form["categories[]"]
			var categories string

			for _, cat := range categ {
				categories = categories + "/" + cat
			}

			err := r.ParseMultipartForm(20 << 20)
			if err != nil {
				AllData.Error = ErrBadSizeImg
				if AllData.ColorMode == "light" {
					err := CreatePostsError.ExecuteTemplate(w, "createposterror.html", AllData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					return
				} else {
					err := DarkCreatePostsError.ExecuteTemplate(w, "createposterror.html", AllData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					return
				}
			}
			var FileName string
			files := r.MultipartForm.File["images[]"]

			if len(files) != 0 {
				for _, fileHeader := range files {
					file, err := fileHeader.Open()
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					defer file.Close()

					if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image/") {
						AllData.Error = ErrBadTypeImg
						if AllData.ColorMode == "light" {
							err := CreatePostsError.ExecuteTemplate(w, "createposterror.html", AllData)
							if err != nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
							}
							return
						} else {
							err := DarkCreatePostsError.ExecuteTemplate(w, "createposterror.html", AllData)
							if err != nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
							}
							return
						}
					}

					fileBytes, err := io.ReadAll(file)
					if err != nil {
						AllData.Error = err
						if AllData.ColorMode == "light" {
							AllData.Error = err
							err := CreatePostsError.ExecuteTemplate(w, "createposterror.html", AllData)
							if err != nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
							}
							return
						} else {
							AllData.Error = err
							err := DarkCreatePostsError.ExecuteTemplate(w, "createposterror.html", AllData)
							if err != nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
							}
							return
						}
					}

					fileName := filepath.Base(fileHeader.Filename)
					destination := "./uploads/" + fileName

					if _, err := os.Stat(destination); err == nil {
						ext := filepath.Ext(fileName)
						fileNameWithoutExt := fileName[:len(fileName)-len(ext)]
						i := 1
						for {
							newFileName := fmt.Sprintf("%s_%d%s", fileNameWithoutExt, i, ext)
							newDestination := filepath.Join("./uploads/", newFileName)
							_, err := os.Stat(newDestination)
							if os.IsNotExist(err) {
								destination = newDestination
								break
							}
							i++
						}
					}

					err = os.WriteFile(destination, fileBytes, 0644)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					FileName = FileName + " " + destination
				}

				errr := CreatePost(Db, AllData.User.User_id, categories, title, content, FileName)
				if errr != nil {
					http.Error(w, errr.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				errr := CreatePost(Db, AllData.User.User_id, categories, title, content, "")
				if errr != nil {
					http.Error(w, errr.Error(), http.StatusInternalServerError)
					return
				}
			}

		} else if postType == "image" {
			imageTitle := r.FormValue("image_title")
			categ := r.Form["categories[]"]
			var categories string

			for _, cat := range categ {
				categories = categories + "/" + cat
			}

			err := r.ParseMultipartForm(20 << 20)
			if err != nil {
				AllData.Error = ErrBadSizeImg
				if AllData.ColorMode == "light" {
					err := CreatePostsError.ExecuteTemplate(w, "createposterror.html", AllData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					return
				} else {
					err := DarkCreatePostsError.ExecuteTemplate(w, "createposterror.html", AllData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					return
				}
			}
			var FileName string
			files := r.MultipartForm.File["images[]"]

			if len(files) != 0 {
				for _, fileHeader := range files {
					file, err := fileHeader.Open()
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					defer file.Close()

					if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image/") {
						AllData.Error = ErrBadTypeImg
						if AllData.ColorMode == "light" {
							err := CreatePostsError.ExecuteTemplate(w, "createposterror.html", AllData)
							if err != nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
							}
							return
						} else {
							err := DarkCreatePostsError.ExecuteTemplate(w, "createposterror.html", AllData)
							if err != nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
							}
							return
						}
					}

					fileBytes, err := io.ReadAll(file)
					if err != nil {
						AllData.Error = err
						if AllData.ColorMode == "light" {
							AllData.Error = err
							err := CreatePostsError.ExecuteTemplate(w, "createposterror.html", AllData)
							if err != nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
							}
							return
						} else {
							AllData.Error = err
							err := DarkCreatePostsError.ExecuteTemplate(w, "createposterror.html", AllData)
							if err != nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
							}
							return
						}
					}

					fileName := filepath.Base(fileHeader.Filename)
					destination := "./uploads/" + fileName

					if _, err := os.Stat(destination); err == nil {
						ext := filepath.Ext(fileName)
						fileNameWithoutExt := fileName[:len(fileName)-len(ext)]
						i := 1
						for {
							newFileName := fmt.Sprintf("%s_%d%s", fileNameWithoutExt, i, ext)
							newDestination := filepath.Join("./uploads/", newFileName)
							_, err := os.Stat(newDestination)
							if os.IsNotExist(err) {
								destination = newDestination
								break
							}
							i++
						}
					}

					err = os.WriteFile(destination, fileBytes, 0644)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					FileName = FileName + " " + destination
				}
			}

			errr := CreatePost(Db, AllData.User.User_id, categories, imageTitle, "", FileName)
			if errr != nil {
				http.Error(w, errr.Error(), http.StatusInternalServerError)
				return
			}
		}
		if err == nil {
			http.Redirect(w, r, "/accueil", http.StatusSeeOther)
			return
		} else {
			if AllData.ColorMode == "light" {
				AllData.Error = err
				err := CreatePostsError.ExecuteTemplate(w, "createposterror.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			} else {
				AllData.Error = err
				err := DarkCreatePostsError.ExecuteTemplate(w, "createposterror.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	}

	if AllData.ColorMode == "light" {
		err = CreatePosts.ExecuteTemplate(w, "createpost.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err = DarkCreatePosts.ExecuteTemplate(w, "createpost.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func PopulairePage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)

	AllData = GetAllDatas(r)

	recherche := r.FormValue("recherche")
	if recherche != "" {
		AllData.AllPosts = GetPostBySearch(Db, recherche, AllData.AllPosts)
	} else {
		AllData.AllPosts, _ = GetAllPostsByLikeCount(Db)
	}

	AllData.AllPosts, _ = GetAllPostsByLikeCount(Db)
	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

	if AllData.User.Email == "" {
		if AllData.ColorMode == "light" {
			err = Populaire.ExecuteTemplate(w, "populaire.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := DarkPopulaire.ExecuteTemplate(w, "populaire.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		AllData.AllNotifications = GetNotifications(Db, UserSession.User_id)
		AllData.AllNotifications = GetNotifications(Db, UserSession.User_id)
		if AllData.ColorMode == "light" {
			err = PopulaireLogged.ExecuteTemplate(w, "populaireLogged.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := DarkPopulaireLogged.ExecuteTemplate(w, "populaireLogged.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

// filtrage a faire !!!!
func PostsPage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)

	AllData = GetAllDatas(r)
	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

	recherche := r.FormValue("recherche")
	if recherche != "" {
		AllData.AllPosts = GetPostBySearch(Db, recherche, AllData.AllPosts)
	} else {
		AllData.AllPosts = GetAllPosts(Db)
	}

	if AllData.User.Email == "" {
		if AllData.ColorMode == "light" {
			err = Posts.ExecuteTemplate(w, "filtragePost.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := DarkPosts.ExecuteTemplate(w, "filtragePost.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		AllData.AllNotifications = GetNotifications(Db, UserSession.User_id)
		if AllData.ColorMode == "light" {
			err = PostsLogged.ExecuteTemplate(w, "filtragePostLogged.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := DarkPostsLogged.ExecuteTemplate(w, "filtragePostLogged.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

// A checker si ça marcher
func NotificationsPage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}

	data, _ := getSessionData(r)
	if data.User.Email == "" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}

	updateUserSession(r)

	AllData = GetAllDatas(r)
	AllData.AllNotifications = GetNotifications(Db, UserSession.User_id)
	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

	recherche := r.FormValue("recherche")
	if recherche != "" {
		AllData.AllNotifications = GetNotifBySearch(Db, UserSession.User_id, recherche)
	} else {
		AllData.AllNotifications = GetNotifications(Db, UserSession.User_id)
	}

	if AllData.ColorMode == "light" {
		err = Notifications.ExecuteTemplate(w, "notification.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err := DarkNotifications.ExecuteTemplate(w, "notification.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ChangeColorMode(w http.ResponseWriter, r *http.Request) {
	updateUserSession(r)
	AllData = GetAllDatas(r)
	if AllData.ColorMode == "light" {
		AllData.ColorMode = "dark"
	} else {
		AllData.ColorMode = "light"
	}

	data, _ := getSessionData(r)

	if UserSession.Role == "user" || UserSession.Role == "admin" {
		createSessionCookie(w, SessionData{
			User: Session{
				UUID:      UserSession.User_id,
				Email:     UserSession.Email,
				Username:  UserSession.Username,
				Role:      UserSession.Role,
				ColorMode: AllData.ColorMode,
			},
		}, 24*time.Hour)
	} else {
		createSessionCookie(w, SessionData{
			User: Session{
				UUID:      data.User.UUID,
				Email:     data.User.Email,
				Username:  data.User.Username,
				Role:      "guest",
				ColorMode: AllData.ColorMode,
			},
		}, 24*time.Hour)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}

	AllData = GetAllDatas(r)

	w.WriteHeader(http.StatusNotFound)

	if AllData.ColorMode == "light" {
		p := "Page not found"
		err = Error404.ExecuteTemplate(w, "error.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		p := "Page not found"
		err = DarkError.ExecuteTemplate(w, "error.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func LikeLogique(w http.ResponseWriter, r *http.Request) {

	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)

	if UserSession.Email == "" {
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	postID := r.URL.RawQuery

	if postID == "" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}

	AllData = GetAllDatas(r)
	postSession, err := GetPost(Db, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if strings.Contains(postSession.Liker, UserSession.Username) {
		UnLikePost(Db, postID, UserSession.Username)
	} else {
		LikePost(Db, postID, UserSession.Username)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func DislikeLogique(w http.ResponseWriter, r *http.Request) {

	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)

	if UserSession.Email == "" {
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	postID := r.URL.RawQuery

	if postID == "" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}

	AllData = GetAllDatas(r)
	postSession, err := GetPost(Db, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if strings.Contains(postSession.Disliker, UserSession.Username) {
		UnDislikePost(Db, postID, UserSession.Username)
	} else {
		DislikePost(Db, postID, UserSession.Username)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func RetweetLogique(w http.ResponseWriter, r *http.Request) {

	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)

	if UserSession.Email == "" {
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	postID := r.URL.RawQuery

	if postID == "" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}

	AllData = GetAllDatas(r)
	postSession, err := GetPost(Db, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if strings.Contains(postSession.Retweeter, UserSession.Username) {
		UnRetweetPost(Db, postID, UserSession.Username)
	} else {
		RetweetPost(Db, postID, UserSession.Username)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func FollowLogique(w http.ResponseWriter, r *http.Request) {

	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)

	if UserSession.Email == "" {
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	user_id := r.URL.RawQuery

	if user_id == "" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}

	AllData = GetAllDatas(r)
	userSession := GetAccount(Db, UserSession.Email)
	userTarget := GetAccountById(Db, user_id)

	if strings.Contains(userSession.FollowingList, userTarget.Username) {
		UpdateUnfollowing(Db, userSession.User_id, userTarget.Username)
		AllData.UserTarget.ImFollowed = false
	} else {
		UpdateFollowing(Db, userSession.User_id, userTarget.Username)
		AllData.UserTarget.ImFollowed = true
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func CategoriePage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}

	if r.URL.RawQuery == "" {
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	updateUserSession(r)

	AllData = GetAllDatas(r)
	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

	categorie_id := r.URL.RawQuery[3:]

	categorie := GetCategoryById(Db, categorie_id)

	recherche := r.FormValue("recherche")
	if recherche != "" {
		AllData.AllPosts = GetPostBySearch(Db, recherche, AllData.AllPosts)
	} else {
		AllData.AllPosts, _ = GetAllPostsByCategorie(Db, categorie.Name)
	}

	AllData.Categorie = categorie

	AllData.AllPosts, _ = GetAllPostsByCategorie(Db, categorie.Name)

	if AllData.User.Email == "" {
		if AllData.ColorMode == "light" {
			err = Categorie.ExecuteTemplate(w, "categorie.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := DarkCategorie.ExecuteTemplate(w, "categorie.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		AllData.AllNotifications = GetNotifications(Db, UserSession.User_id)
		if AllData.ColorMode == "light" {
			err = CategorieLogged.ExecuteTemplate(w, "categorieLogged.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := DarkCategorieLogged.ExecuteTemplate(w, "categorieLogged.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func ReportLogique(w http.ResponseWriter, r *http.Request) {
	updateUserSession(r)

	if UserSession.Email == "" {
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	postID := r.URL.RawQuery

	if postID == "" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}

	ReportPost(Db, postID)

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

// Page une fois après avoir fais la demande de changement de mot de passe (reussite)
func ForgotPasswordSuccessPage(w http.ResponseWriter, r *http.Request) {
	AllData = GetAllDatas(r)

	if AllData.ColorMode == "light" {
		err := ForgotPasswordSuccess.ExecuteTemplate(w, "password_reset_success.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err := DarkForgotPasswordSuccess.ExecuteTemplate(w, "password_reset_success.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Page pour écrire le mail pour changer de mot de passe
func ForgotPasswordPage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)

	AllData = GetAllDatas(r)

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		if email == "" {
			if AllData.ColorMode == "light" {
				err := ForgotPasswordError.ExecuteTemplate(w, "forgotpassworderror.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			} else {
				err := DarkForgotPasswordError.ExecuteTemplate(w, "forgotpassworderror.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
		if FindAccount(Db, email) {
			token := EncodeToken(email)

			CreateRoute(w, r, token)
			SendPasswordResetEmail(email, token)
			if email != "" {
				http.Redirect(w, r, "/forgot-password-success", http.StatusSeeOther)
				return
			}
		} else {
			http.Redirect(w, r, "/no-mail-found", http.StatusSeeOther)
			return
		}
	}

	if AllData.ColorMode == "light" {
		err = ForgotPassword.ExecuteTemplate(w, "forgotpassword.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err = DarkForgotPassword.ExecuteTemplate(w, "forgotpassword.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Page pour écrire le nouveau mdp si le token est valide
func ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(r)
	AllData = GetAllDatas(r)

	token := r.FormValue("token")
	if ResetPasswordMap[token] == "valid" {
		if token == "" {
			http.Error(w, "Token manquant dans l'URL", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPost {
			newPassword := r.FormValue("password")
			checkPassword := r.FormValue("checkpassword")

			// condition du mot de passe
			if newPassword != "" && newPassword == checkPassword {
				token = strings.Trim(token, "/")

				email, _ := DecodeToken(token)
				user := GetAccount(Db, email)

				ChangePassword(Db, user.User_id, newPassword)
				http.Redirect(w, r, "/reset-password-success", http.StatusSeeOther)
				return
			} else {
				if AllData.ColorMode == "light" {
					err := ResetPasswordError.ExecuteTemplate(w, "reset_password_error.html", AllData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					return
				} else {
					err := DarkResetPasswordError.ExecuteTemplate(w, "reset_password_error.html", AllData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					return
				}
			}
		}
		if AllData.ColorMode == "light" {
			err := ResetPassword.ExecuteTemplate(w, "reset_password.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := DarkResetPassword.ExecuteTemplate(w, "reset_password.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		http.Redirect(w, r, "/lien-expire", http.StatusSeeOther)
	}
}

// Page lorsque le mdp est modifié
func PasswordResetSuccessPage(w http.ResponseWriter, r *http.Request) {
	updateUserSession(r)
	AllData = GetAllDatas(r)
	ResetPasswordMap[URL] = "invalid"
	if AllData.ColorMode == "light" {
		err := PasswordResetSuccess.ExecuteTemplate(w, "reset_success.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err := DarkPasswordResetSuccess.ExecuteTemplate(w, "reset_success.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Page lorsque l'email n'est pas trouvé
func NoMailFoundPage(w http.ResponseWriter, r *http.Request) {
	updateUserSession(r)
	AllData = GetAllDatas(r)

	if AllData.ColorMode == "light" {
		err := NoMailFound.ExecuteTemplate(w, "noMailFound.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err := DarkNoMailFound.ExecuteTemplate(w, "noMailFound.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Page lorsque le lien est expiré
func ExpiredLinkPage(w http.ResponseWriter, r *http.Request) {
	updateUserSession(r)
	AllData = GetAllDatas(r)

	if AllData.ColorMode == "light" {
		err := ExpiredLink.ExecuteTemplate(w, "expiredLink.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err := DarkExpiredLink.ExecuteTemplate(w, "expiredLink.html", AllData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ReglagePage(w http.ResponseWriter, r *http.Request) {
	updateUserSession(r)
	AllData = GetAllDatas(r)
	AllData.AllNotifications = GetNotifications(Db, UserSession.User_id)
	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

	settings := r.URL.RawQuery

	if settings == "profile" {
		if UserSession.Email == "" {
			http.Redirect(w, r, "/connexion", http.StatusSeeOther)
			return
		} else {
			if AllData.ColorMode == "light" {
				err := ReglageProfile.ExecuteTemplate(w, "reglageVotreProfile.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			} else {
				err := DarkReglageProfile.ExecuteTemplate(w, "reglageVotreProfile.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	} else if settings == "prenium" {
		if UserSession.Email == "" {
			http.Redirect(w, r, "/connexion", http.StatusSeeOther)
			return
		} else {
			if AllData.ColorMode == "light" {
				err := ReglagePrenium.ExecuteTemplate(w, "prenium.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			} else {
				err := DarkReglagePrenium.ExecuteTemplate(w, "prenium.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	} else if settings == "assist" {
		if AllData.ColorMode == "light" {
			err := ReglageAssist.ExecuteTemplate(w, "assist.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		} else {
			err := DarkReglageAssist.ExecuteTemplate(w, "assist.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	} else if settings == "profile/account" {
		if UserSession.Email == "" {
			http.Redirect(w, r, "/connexion", http.StatusSeeOther)
			return
		} else {
			if AllData.ColorMode == "light" {
				err := ReglageInfo.ExecuteTemplate(w, "reglageInfoCompte.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			} else {
				err := DarkReglageInfo.ExecuteTemplate(w, "reglageInfoCompte.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	} else if settings == "profile/change-password" {
		if UserSession.Email == "" {
			http.Redirect(w, r, "/connexion", http.StatusSeeOther)
			return
		} else {
			if AllData.ColorMode == "light" {
				err := ReglageChangePassword.ExecuteTemplate(w, "reglageChangerDeMotDePasse.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			} else {
				err := DarkReglageChangePassword.ExecuteTemplate(w, "reglageChangerDeMotDePasse.html", AllData)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	}

	if UserSession.Email == "" {
		if AllData.ColorMode == "light" {
			err := Reglage.ExecuteTemplate(w, "reglage.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := DarkReglage.ExecuteTemplate(w, "reglage.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {

		if AllData.ColorMode == "light" {
			err := ReglageLogged.ExecuteTemplate(w, "reglageLogged.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := DarkReglageLogged.ExecuteTemplate(w, "reglageLogged.html", AllData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
