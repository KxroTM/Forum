package forum

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	User_id      string
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
var AllUsers []User

func CreateUser(db *sql.DB, username, email, password string) {
	currentTime := time.Now()
	time := currentTime.Format("02-01-2006")

	u, err := uuid.NewV4()
	if err != nil {
		fmt.Println("Erreur lors de la génération de l'UUID :", err)
		return
	}

	userSession = User{
		User_id:      u.String(),
		Role:         "user",
		Username:     username,
		Email:        email,
		Password:     hashPasswordSHA256(password),
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

	AllUsers = append(AllUsers, userSession)
	UpdateDb(db)
}

func DeleteUser(db *sql.DB, user_id int) {
	db.Exec(`DELETE FROM users WHERE UUID = ?`, user_id)
}

// Function for hash a password
func hashPasswordSHA256(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

// Return true if the username is in the DB, false if not
func findAccount(user_id string) bool {
	for _, user := range AllUsers {
		if user.User_id == user_id {
			return true
		}
	}
	return false
}

func UpdateDb(db *sql.DB) {
	// Sélectionner toutes les lignes de la table User
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Printf("erreur lors de la récupération des utilisateurs depuis la base de données: %v", err)
	}
	defer rows.Close()

	// Initialiser une slice pour stocker les utilisateurs
	var users []User

	// Parcourir les lignes résultantes
	for rows.Next() {
		var user User
		// Scanner les valeurs des colonnes dans la structure User
		if err := rows.Scan(&user.User_id, &user.Role, &user.Username, &user.Email, &user.Password, &user.CreationDate, &user.UpdateDate, &user.Pfp, &user.Bio, &user.Links, &user.CategorieSub, &user.Follower, &user.Following); err != nil {
			fmt.Printf("erreur lors de la lecture des données utilisateur depuis la base de données: %v", err)
		}
		// Ajouter l'utilisateur à la slice
		users = append(users, user)
	}
	// Vérifier s'il y a eu une erreur lors de l'itération des lignes
	if err := rows.Err(); err != nil {
		fmt.Printf("erreur lors de la récupération des utilisateurs depuis la base de données: %v", err)
	}

	// Mettre à jour la structure AllUsers avec les utilisateurs récupérés
	AllUsers = users
}

func DeleteAllUsers(db *sql.DB) {
	db.Exec("DELETE FROM users")
}
