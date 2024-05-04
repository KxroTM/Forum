package forum

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/gofrs/uuid"
)

var Db *sql.DB

// TYPE OF USER :
// - guests
// - users
// - prenium users
// - moderators
// - administrators

type User struct {
	User_id       string
	Role          string
	Username      string
	Email         string
	Password      string
	CreationDate  string
	UpdateDate    string
	Pfp           string
	Bio           string
	Links         string
	CategorieSub  string
	Follower      int
	FollowerList  string
	Following     int
	FollowingList string
}

var UserSession User
var AllUsers []User

var banWords = []string{
	"idiot", "imbécile", "crétin", "con", "abruti", "connard", "enfoiré", "salopard",
	"sexe", "porno", "XXX", "nue", "seins", "cul", "bite", "vagin", "pénis", "orgasme", "éjaculation",
	"raciste", "homophobe", "sexiste", "islamophobe", "antisémite", "xénophobe", "suprémaciste", "haineux",
	"drogue", "vol", "fraude", "escroc", "trafic", "prostitué", "pédophile", "viol", "meurtre", "terroriste",
	"Dieu", "Allah", "Jésus", "Satan", "Lucifer", "messie", "prophète", "pape", "imam", "rabbin",
	"merde", "baise", "putain", "foutre", "enculé", "niquer", "chier", "salaud", "cul", "bite",
	"tuer", "battre", "torture", "maltraitance", "séquestration", "cruauté", "violence", "massacre",
	"haine", "violence", "assassinat", "extermination", "guerre", "destruction", "attaquer", "détruire", "anéantir",
	"nazisme", "communisme", "fascisme", "dictature", "totalitarisme", "extrémisme", "nationalisme", "anarchie",
	"Trump", "Hitler", "Staline", "Mao", "Ben Laden", "Saddam Hussein", "Kim Jong-un", "Poutine", "Assad",
}

func UpdateUserDb(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.User_id, &user.Role, &user.Username, &user.Email, &user.Password, &user.CreationDate, &user.UpdateDate, &user.Pfp, &user.Follower, &user.Following, &user.Bio, &user.Links, &user.CategorieSub, &user.FollowerList, &user.FollowingList)
		if err != nil {
			return err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	AllUsers = users
	return nil
}

func SignUpUser(db *sql.DB, username, email, password string) error {
	currentTime := time.Now()
	time := currentTime.Format("02-01-2006")

	u, err := uuid.NewV4()
	if err != nil {
		return err
	}

	UserSession = User{
		User_id:       u.String(),
		Role:          "user",
		Username:      username,
		Email:         email,
		Password:      hashPasswordSHA256(password),
		CreationDate:  time,
		UpdateDate:    time,
		Pfp:           "",
		Bio:           "",
		Links:         "",
		CategorieSub:  "",
		Follower:      0,
		FollowerList:  "",
		Following:     0,
		FollowingList: "",
	}

	pfp := rand.Intn(4) + 1
	switch pfp {
	case 1:
		UserSession.Pfp = "../../style/media/default_avatar/avatar_01.png"
	case 2:
		UserSession.Pfp = "../../style/media/default_avatar/avatar_02.png"
	case 3:
		UserSession.Pfp = "../../style/media/default_avatar/avatar_03.png"
	case 4:
		UserSession.Pfp = "../../style/media/default_avatar/avatar_04.png"
	}

	db.Exec(`INSERT INTO users (UUID, role, username, email, password, created_at, updated_at, profilePicture, followers, following, bio, links, categoriesSub, followersList, followingList) 
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, UserSession.User_id, UserSession.Role, UserSession.Username, UserSession.Email, UserSession.Password, UserSession.CreationDate, UserSession.UpdateDate, UserSession.Pfp, UserSession.Follower, UserSession.Following, UserSession.Bio, UserSession.Links, UserSession.CategorieSub, UserSession.FollowerList, UserSession.FollowingList)

	AllUsers = append(AllUsers, UserSession)
	UpdateUserDb(db)
	return nil
}

func LoginUser(db *sql.DB, email, password string) bool {
	for _, user := range AllUsers {
		if user.Email == email && user.Password == hashPasswordSHA256(password) {
			UserSession = user
			return true
		}
	}
	return false
}

func LogoutUser() {
	UserSession = User{
		Role: "guest",
	}
}

func DeleteUser(db *sql.DB, user_id string) error {
	db.Exec(`DELETE FROM users WHERE UUID = ?`, user_id)
	err := UpdateUserDb(db)
	if err != nil {
		return err
	}
	return nil
}

func hashPasswordSHA256(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func FindAccount(email string) bool {
	for _, user := range AllUsers {
		if user.Email == email {
			return true
		}
	}
	return false
}

func GetAccount(email string) User {
	for _, user := range AllUsers {
		if user.Email == email {
			return user
		}
	}
	return User{}
}

func GetAccountById(user_id string) User {
	for _, user := range AllUsers {
		if user.User_id == user_id {
			return user
		}
	}
	return User{}
}

func GetAccountByUsername(username string) User {
	for _, user := range AllUsers {
		if user.Username == username {
			return user
		}
	}
	return User{}
}

func DeleteAllUsers(db *sql.DB) error {
	db.Exec("DELETE FROM users")
	err := UpdateUserDb(db)
	if err != nil {
		return err
	}
	return nil
}

func ChangePassword(db *sql.DB, user_id string, newPassword string) error {
	db.Exec(`UPDATE users SET password = ? WHERE UUID = ?`, hashPasswordSHA256(newPassword), user_id)
	err := UpdateUserDb(db)
	if err != nil {
		return err
	}
	return nil
}

func IsPasswordValid(password string) bool {
	if len(password) < 8 {
		return false
	}

	lowerCheck := false
	upperCheck := false
	digitCheck := false
	specialCharCheck := false

	specialCharsPattern := `[^a-zA-Z0-9]`
	regexp := regexp.MustCompile(specialCharsPattern)

	for _, char := range password {
		if unicode.IsLower(char) {
			lowerCheck = true
		}
		if unicode.IsUpper(char) {
			upperCheck = true
		}
		if unicode.IsDigit(char) {
			digitCheck = true
		}
		if regexp.MatchString(string(char)) {
			specialCharCheck = true
		}
	}

	return lowerCheck && upperCheck && digitCheck && specialCharCheck
}

func ChangeUsername(db *sql.DB, user_id string, newUsername string) error {
	db.Exec(`UPDATE users SET username = ? WHERE UUID = ?`, newUsername, user_id)
	err := UpdateUserDb(db)
	if err != nil {
		return err
	}
	return nil
}

func isUsernameAvailable(username string) bool {
	for _, user := range AllUsers {
		if user.Username == username {
			return false
		}
	}
	return true
}

func IsUsernameValid(username string) bool {
	if (len(username) > 4 || len(username) < 15) && !containsBanWord(strings.ToLower(username)) {
		return isUsernameAvailable(username)
	}
	return false
}

func isEmailAvailable(email string) bool {
	for _, user := range AllUsers {
		if user.Email == email {
			return false
		}
	}
	return true
}

func IsEmailValid(email string) bool {
	emailPatern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	regex := regexp.MustCompile(emailPatern)

	return regex.MatchString(email) && isEmailAvailable(email)
}

func SetModerator(db *sql.DB, user_id string) error {
	db.Exec(`UPDATE users SET role = ? WHERE UUID = ?`, "moderator", user_id)
	err := UpdateUserDb(db)
	if err != nil {
		return err
	}
	return nil
}

func containsBanWord(word string) bool {
	for _, words := range banWords {
		if strings.EqualFold(strings.ToLower(removeAccents(words)), removeAccents(word)) {
			return true
		}
	}
	return false
}

func removeAccents(s string) string {
	t := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.Is(unicode.Mn, r) {
			t = append(t, r)
		}
	}
	return string(t)
}

func UpdateDate(db *sql.DB, user_id string) error {
	currentTime := time.Now()
	time := currentTime.Format("02-01-2006")
	db.Exec(`UPDATE users SET updated_at = ? WHERE UUID = ?`, time, user_id)
	err := UpdateUserDb(db)
	if err != nil {
		return err
	}
	return nil
}

func UpdateProfilePicture(db *sql.DB, user_id string, pfp string) (bool, error) {
	if UserSession.Role == "user" && isProfilePictureNotAGif(pfp) {
		db.Exec(`UPDATE users SET profilePicture = ? WHERE UUID = ?`, pfp, user_id)
		err := UpdateUserDb(db)
		if err != nil {
			return false, err
		}
		return true, nil
	} else if UserSession.Role != "user" {
		db.Exec(`UPDATE users SET profilePicture = ? WHERE UUID = ?`, pfp, user_id)
		err := UpdateUserDb(db)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func isProfilePictureNotAGif(pfp string) bool {
	return !strings.HasSuffix(strings.ToLower(pfp), ".gif") && !strings.HasSuffix(strings.ToLower(pfp), ".apng")
}

func UpdateBio(db *sql.DB, user_id string, bio string) error {
	db.Exec(`UPDATE users SET bio = ? WHERE UUID = ?`, bio, user_id)
	err := UpdateUserDb(db)
	if err != nil {
		return err
	}
	return nil
}

func UpdateLinks(db *sql.DB, user_id string, links string) error {
	db.Exec(`UPDATE users SET links = ? WHERE UUID = ?`, links, user_id)
	err := UpdateUserDb(db)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCategoriesSub(db *sql.DB, user_id string, categorie string) error {
	categoriesSub := UserSession.CategorieSub + "," + categorie
	db.Exec(`UPDATE users SET categoriesSub = ? WHERE UUID = ?`, categoriesSub, user_id)
	err := UpdateUserDb(db)
	if err != nil {
		return err
	}
	return nil
}

func UpdateFollowing(db *sql.DB, user_id string, username string) error { // username etant la personne que l'on va follow

	userToFollow := GetAccount(username)

	// Mise a jour de notre nombre de following
	db.Exec(`UPDATE users SET following = ? WHERE UUID = ?`, UserSession.Following+1, user_id)

	// Mise a jour de la liste des personnes que l'on follow
	db.Exec(`UPDATE users SET followingList = ? WHERE UUID = ?`, UserSession.FollowingList+","+username, user_id)

	// Mise a jour du nombre de followers de la personne que l'on follow
	db.Exec(`UPDATE users SET followers = ? WHERE UUID = ?`, userToFollow.Follower+1, userToFollow.User_id)

	// Mise a jour de la liste des followers de la personne que l'on follow
	db.Exec(`UPDATE users SET followersList = ? WHERE UUID = ?`, userToFollow.FollowerList+","+UserSession.Username, userToFollow.User_id)

	err := UpdateUserDb(db)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUnfollowing(db *sql.DB, user_id string, username string) error { // username etant la personne que l'on va unfollow

	userToUnfollow := GetAccount(username)

	// Mise a jour de notre nombre de following
	db.Exec(`UPDATE users SET following = ? WHERE UUID = ?`, UserSession.Following-1, user_id)

	// Mise a jour de la liste des personnes que l'on unfollow
	db.Exec(`UPDATE users SET followingList = ? WHERE UUID = ?`, strings.Replace(UserSession.FollowingList, ","+username, "", -1), user_id)

	// Mise a jour du nombre de followers de la personne que l'on unfollow
	db.Exec(`UPDATE users SET followers = ? WHERE UUID = ?`, userToUnfollow.Follower-1, userToUnfollow.User_id)

	// Mise a jour de la liste des followers de la personne que l'on unfollow
	db.Exec(`UPDATE users SET followersList = ? WHERE UUID = ?`, strings.Replace(userToUnfollow.FollowerList, ","+UserSession.Username, "", -1), userToUnfollow.User_id)

	err := UpdateUserDb(db)
	if err != nil {
		return err
	}
	return nil
}

func GetAllMail() []string {
	var mails []string
	for _, user := range AllUsers {
		mails = append(mails, user.Email)
	}
	return mails
}

func GetAllUser() []User {
	return AllUsers
}

func GetAllDatas() DataStruct {
	return DataStruct{
		User:       UserSession,
		UserTarget: User{},
		AllUsers:   GetAllUser(),
		Post:       Post{},
		AllPosts:   GetAllPosts(),
		Comment:    Comment{},
		// AllComments:      GetAllComments(),
		Notification: Notification{},
		// AllNotifications: GetAllNotifications(),
	}
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
