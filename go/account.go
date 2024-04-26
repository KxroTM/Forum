package forum

import (
	"database/sql"
	"time"
)

type User struct {
	User_id      int
	Role         string
	Username     string
	Email        string
	Password     string
	CreationDate string
	UpdateDate   string
	Pfp          string
	Bio          string
	Links        string
	CategorieSub string
	Follower     int
	Following    int
}

var userSession User

var ID = 0

func CreateUser(db *sql.DB, username, email, password string) {
	currentTime := time.Now()
	time := currentTime.Format("02-01-2006")

	userSession = User{
		User_id:      ID + 1,
		Role:         "user",
		Username:     username,
		Email:        email,
		Password:     password,
		CreationDate: time,
		UpdateDate:   time,
		Pfp:          "",
		Bio:          "",
		Links:        "",
		CategorieSub: "",
		Follower:     0,
		Following:    0,
	}

	db.Exec(`INSERT INTO users (UUID, role, username, email, password, created_at, updated_at, profilePicture, followers, following, bio, links, categoriesSub) 
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, userSession.User_id, userSession.Role, userSession.Username, userSession.Email, userSession.Password, userSession.CreationDate, userSession.UpdateDate, userSession.Pfp, userSession.Follower, userSession.Following, userSession.Bio, userSession.Links, userSession.CategorieSub)
}

func DeleteUser(db *sql.DB, user_id int) {
	db.Exec(`DELETE FROM users WHERE UUID = ?`, user_id)
}
