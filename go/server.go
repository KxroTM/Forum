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
	updateUserSession(w, r)

	AllData = GetAllDatas(w, r)

	// Si l'utilisateur est déjà connecté, on le redirige vers la page d'accueil
	data, _ := getSessionData(w, r)
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
	updateUserSession(w, r)

	AllData = GetAllDatas(w, r)

	// Si l'utilisateur est déjà connecté, on le redirige vers la page d'accueil
	data, _ := getSessionData(w, r)
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
	updateUserSession(w, r)

	query := r.URL.RawQuery
	AllData = GetAllDatas(w, r)

	if query == "pourtoi" {
		AllData.AllPosts = ForYouPageAlgorithm(Db, UserSession.User_id)
	} else if query == "suivies" {
		AllData.AllPosts = GetPostByFollowing(Db, UserSession.User_id)
	}

	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

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
}

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(w, r)
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

	AllData = GetAllDatas(w, r)
	AllData.UserTarget = GetAccountByUsername(Db, username)
	AllData.AllPosts, _ = GetAllPostsByUser(Db, AllData.UserTarget.User_id)
	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

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
}

func PostPage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(w, r)
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

	AllData = GetAllDatas(w, r)
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
}

func CreatePostPage(w http.ResponseWriter, r *http.Request) {
	data, _ := getSessionData(w, r)
	if data.User.Email == "" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}

	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(w, r)

	AllData = GetAllDatas(w, r)
	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

	if r.Method == http.MethodPost {
		postType := r.FormValue("post_type")
		if postType == "text" {
			title := r.FormValue("title")
			content := r.FormValue("text")
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

				errr := CreatePost(Db, AllData.User.User_id, "", title, content, FileName)
				if errr != nil {
					http.Error(w, errr.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				errr := CreatePost(Db, AllData.User.User_id, "", title, content, "")
				if errr != nil {
					http.Error(w, errr.Error(), http.StatusInternalServerError)
					return
				}
			}

		} else if postType == "image" {
			imageTitle := r.FormValue("image_title")
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

			errr := CreatePost(Db, AllData.User.User_id, "", imageTitle, "", FileName)
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
	updateUserSession(w, r)

	AllData = GetAllDatas(w, r)
	AllData.AllPosts, _ = GetAllPostsByLikeCount(Db)
	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

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
}

func PostsPage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(w, r)

	AllData = GetAllDatas(w, r)
	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

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
}

func NotificationsPage(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	err := IPsLog(clientIP + "  ==>  " + r.URL.Path)
	if err != nil {
		log.Println(err)
	}
	updateUserSession(w, r)

	AllData = GetAllDatas(w, r)
	// AllData.AllNotifications, _ = GetNotifications(Db, UserSession.User_id)
	AllData.RecommendedUser = RecommendedUsers(Db, UserSession.User_id)

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
	updateUserSession(w, r)
	AllData = GetAllDatas(w, r)
	if AllData.ColorMode == "light" {
		AllData.ColorMode = "dark"
	} else {
		AllData.ColorMode = "light"
	}

	data, _ := getSessionData(w, r)

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

	AllData = GetAllDatas(w, r)

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
