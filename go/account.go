package forum

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"net/http"
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
// - premium users
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
	IsMyAccount   bool
	ImFollowed    bool
	HeFollowed    bool
}

var UserSession User

var banWords = []string{
	"idiot", "imbecile", "cretin", "con", "abruti", "connard", "enfoire", "salopard",
	"sexe", "porno", "XXX", "nue", "seins", "cul", "bite", "vagin", "penis", "orgasme", "ejaculation",
	"raciste", "homophobe", "sexiste", "islamophobe", "antisemite", "xenophobe", "supremaciste", "haineux",
	"drogue", "vol", "fraude", "escroc", "trafic", "prostitue", "pedophile", "viol", "meurtre", "terroriste",
	"dieu", "jesus", "satan", "lucifer", "messie", "prophete", "pape", "imam", "rabbin",
	"merde", "baise", "putain", "foutre", "encule", "niquer", "chier", "salaud", "cul", "bite",
	"tuer", "battre", "torture", "maltraitance", "sequestration", "cruaute", "violence", "massacre",
	"haine", "violence", "assassinat", "extermination", "guerre", "destruction", "attaquer", "detruire", "aneantir",
	"nazisme", "communisme", "fascisme", "dictature", "totalitarisme", "extremisme", "nationalisme", "anarchie",
	"trump", "hitler", "staline", "mao", "benladen", "saddamhussein", "laden", "hussein", "kimjong-un", "poutine", "assad", "fdp", "arabe",
	"nazi", "youssef", "youss", "yous", "gay", "pd", "lgbt", "homo", "bz", "ntm", "tamere", "mere", "nique", "tue", "extermine", "israel", "kippa",
	"noir", "negre", "negro", "singe", "bardella",
}

func checkAllConditionsSignUp(db *sql.DB, username, email, password, passwordcheck string) error {
	if strings.Contains(username, " ") || strings.Contains(username, "-") {
		return ErrSpaceInUsername
	}
	if strings.Contains(password, " ") {
		return ErrSpaceInPassword
	}
	if username == "" {
		return ErrEmptyFieldPseudo
	}

	if email == "" {
		return ErrEmptyFieldEmail
	}

	if password == "" {
		return ErrEmptyFieldPassword
	}

	if passwordcheck == "" {
		return ErrEmptyFieldPasswordCheck
	}

	if password != passwordcheck {
		return ErrInvalidPasswordCheck
	}

	if !IsUsernameValid(username) {
		return ErrInvalidPseudo
	}

	if !IsEmailValid(email) {
		return ErrInvalidEmail
	}

	if !IsPasswordValid(password) {
		return ErrInvalidPassword
	}

	if !isEmailAvailable(db, email) {
		return ErrMailAlreadyUsed
	}

	if !isUsernameAvailable(db, username) {
		return ErrPseudoAlreadyUsed
	}
	return nil
}

func SignUpUser(db *sql.DB, username, email, password, passwordcheck string) error {

	err := checkAllConditionsSignUp(db, username, email, password, passwordcheck)
	if err != nil {
		return err
	}

	time := time.Now().Format("02-01-2006")

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

	SendCreatedAccountEmail(email, username)
	return nil
}

func SignUpUserOauth(db *sql.DB, username, email, password string) error {

	time := time.Now().Format("02-01-2006")

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

	SendCreatedAccountEmail(email, username)
	return nil
}

func LoginUser(db *sql.DB, email, password string) (bool, error) {
	if email == "" {
		return false, ErrEmptyFieldEmail
	}

	if password == "" {
		return false, ErrEmptyFieldPassword
	}

	query := "SELECT UUID, role, username, email, password, created_at, updated_at, profilePicture, bio, links, categoriesSub, followers, followersList, following, followingList FROM users WHERE email = ?"
	row := db.QueryRow(query, email)

	var user User
	err := row.Scan(&user.User_id, &user.Role, &user.Username, &user.Email, &user.Password, &user.CreationDate, &user.UpdateDate, &user.Pfp, &user.Bio, &user.Links, &user.CategorieSub, &user.Follower, &user.FollowerList, &user.Following, &user.FollowingList)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, ErrBadEmail
		}
		return false, err
	}

	if user.Password != hashPasswordSHA256(password) {
		return false, ErrBadPassword
	}

	UserSession = user

	return true, nil
}

func LogoutUser() {
	UserSession = User{
		Role: "guest",
	}
}

func DeleteUser(db *sql.DB, user_id string) {
	db.Exec(`DELETE FROM users WHERE UUID = ?`, user_id)
}

func hashPasswordSHA256(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func FindAccount(db *sql.DB, email string) bool {
	query := "SELECT COUNT(*) FROM users WHERE email = ?"
	var count int
	err := db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func GetAccount(db *sql.DB, email string) User {
	query := "SELECT UUID, role, username, email, password, created_at, updated_at, profilePicture, bio, links, categoriesSub, followers, followersList, following, followingList FROM users WHERE email = ?"
	var user User
	err := db.QueryRow(query, email).Scan(&user.User_id, &user.Role, &user.Username, &user.Email, &user.Password, &user.CreationDate, &user.UpdateDate, &user.Pfp, &user.Bio, &user.Links, &user.CategorieSub, &user.Follower, &user.FollowerList, &user.Following, &user.FollowingList)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}
		}
		return User{}
	}
	return user
}

func GetAccountById(db *sql.DB, user_id string) User {

	query := "SELECT UUID, role, username, email, password, created_at, updated_at, profilePicture, bio, links, categoriesSub, followers, followersList, following, followingList FROM users WHERE UUID = ?"
	var user User
	err := db.QueryRow(query, user_id).Scan(&user.User_id, &user.Role, &user.Username, &user.Email, &user.Password, &user.CreationDate, &user.UpdateDate, &user.Pfp, &user.Bio, &user.Links, &user.CategorieSub, &user.Follower, &user.FollowerList, &user.Following, &user.FollowingList)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}
		}
		return User{}
	}
	return user
}

func GetAccountByUsername(db *sql.DB, username string) User {
	query := "SELECT UUID, role, username, email, password, created_at, updated_at, profilePicture, bio, links, categoriesSub, followers, followersList, following, followingList FROM users WHERE username = ?"
	var user User
	err := db.QueryRow(query, username).Scan(&user.User_id, &user.Role, &user.Username, &user.Email, &user.Password, &user.CreationDate, &user.UpdateDate, &user.Pfp, &user.Bio, &user.Links, &user.CategorieSub, &user.Follower, &user.FollowerList, &user.Following, &user.FollowingList)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}
		}
		return User{}
	}

	if UserSession.Username == user.Username {
		user.IsMyAccount = true
	}

	if strings.Contains(UserSession.FollowingList, user.Username) {
		user.ImFollowed = true
	}

	if strings.Contains(UserSession.FollowerList, user.Username) {
		user.HeFollowed = true
	}

	return user
}

func DeleteAllUsers(db *sql.DB) {
	db.Exec("DELETE FROM users")
}

func ChangePassword(db *sql.DB, user_id string, newPassword string) {
	time := time.Now().Format("02-01-2006")
	db.Exec(`UPDATE users SET password = ?, updated_at = ? WHERE UUID = ?`, hashPasswordSHA256(newPassword), time, user_id)
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

func ChangeUsername(db *sql.DB, user_id string, newUsername string) {
	db.Exec(`UPDATE users SET username = ? WHERE UUID = ?`, newUsername, user_id)
}

func isUsernameAvailable(db *sql.DB, username string) bool {
	query := "SELECT COUNT(*) FROM users WHERE username = ?"
	var count int
	err := db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false
	}
	return count == 0
}

func IsUsernameValid(username string) bool {
	if (len(username) > 4 && len(username) < 15) && !containsBanWord(username) {
		return isUsernameAvailable(Db, username)
	}
	return false
}

func isEmailAvailable(db *sql.DB, email string) bool {
	query := "SELECT COUNT(*) FROM users WHERE email = ?"
	var count int
	err := db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false
	}
	return count == 0
}

func IsEmailValid(email string) bool {
	emailPatern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	regex := regexp.MustCompile(emailPatern)

	return regex.MatchString(email)
}

func SetModerator(db *sql.DB, user_id string) {
	db.Exec(`UPDATE users SET role = ? WHERE UUID = ?`, "moderator", user_id)
}

func containsBanWord(word string) bool {
	word = strings.ToLower(removeAccents(word))
	for _, banWord := range banWords {
		if strings.Contains(word, banWord) {
			return true
		}
	}
	return false
}

func removeAccents(s string) string {
	t := make([]rune, len(s))
	for i, r := range s {
		switch r {
		case 'À', 'Á', 'Â', 'Ã', 'Ä', 'Å':
			t[i] = 'A'
		case 'à', 'á', 'â', 'ã', 'ä', 'å':
			t[i] = 'a'
		case 'Ç':
			t[i] = 'C'
		case 'ç':
			t[i] = 'c'
		case 'È', 'É', 'Ê', 'Ë':
			t[i] = 'E'
		case 'è', 'é', 'ê', 'ë':
			t[i] = 'e'
		case 'Î', 'Ï', 'Í', 'Ì':
			t[i] = 'I'
		case 'î', 'ï', 'í', 'ì':
			t[i] = 'i'
		case 'Ñ':
			t[i] = 'N'
		case 'ñ':
			t[i] = 'n'
		case 'Ò', 'Ó', 'Ô', 'Õ', 'Ö':
			t[i] = 'O'
		case 'ò', 'ó', 'ô', 'õ', 'ö':
			t[i] = 'o'
		case 'Ù', 'Ú', 'Û', 'Ü':
			t[i] = 'U'
		case 'ù', 'ú', 'û', 'ü':
			t[i] = 'u'
		case 'Ý', 'Ỳ', 'Ỹ', 'Ỷ', 'Ỵ':
			t[i] = 'Y'
		case 'ý', 'ỳ', 'ỹ', 'ỷ', 'ỵ':
			t[i] = 'y'
		default:
			t[i] = r
		}
	}
	return string(t)
}

func UpdateDate(db *sql.DB, user_id string) {
	currentTime := time.Now()
	time := currentTime.Format("02-01-2006")
	db.Exec(`UPDATE users SET updated_at = ? WHERE UUID = ?`, time, user_id)
}

func UpdateProfilePicture(db *sql.DB, user_id string, pfp string) bool {
	if UserSession.Role == "user" && isProfilePictureNotAGif(pfp) {
		db.Exec(`UPDATE users SET profilePicture = ? WHERE UUID = ?`, pfp, user_id)
		return true
	} else if UserSession.Role != "user" {
		db.Exec(`UPDATE users SET profilePicture = ? WHERE UUID = ?`, pfp, user_id)
		return true
	}
	return false
}

func isProfilePictureNotAGif(pfp string) bool {
	return !strings.HasSuffix(strings.ToLower(pfp), ".gif") && !strings.HasSuffix(strings.ToLower(pfp), ".apng")
}

func UpdateBio(db *sql.DB, user_id string, bio string) {
	db.Exec(`UPDATE users SET bio = ? WHERE UUID = ?`, bio, user_id)
}

func UpdateLinks(db *sql.DB, user_id string, links string) {
	db.Exec(`UPDATE users SET links = ? WHERE UUID = ?`, links, user_id)
}

func AddCategoriesSub(db *sql.DB, user_id string, categorie string) {
	categoriesSub := UserSession.CategorieSub + "," + categorie
	db.Exec(`UPDATE users SET categoriesSub = ? WHERE UUID = ?`, categoriesSub, user_id)
	db.Exec(`UPDATE categories SET users = users + 1 WHERE name = ?`, categorie)
}

func RemoveCategoriesSub(db *sql.DB, user_id string, categorie string) {
	categoriesSub := strings.Replace(UserSession.CategorieSub, ","+categorie, "", -1)
	db.Exec(`UPDATE users SET categoriesSub = ? WHERE UUID = ?`, categoriesSub, user_id)
	db.Exec(`UPDATE categories SET users = users - 1 WHERE name = ?`, categorie)
}

func UpdateFollowing(db *sql.DB, user_id, username string) { // username etant la personne que l'on va follow

	userToFollow := GetAccountByUsername(db, username)

	// Mise a jour de notre nombre de following
	db.Exec(`UPDATE users SET following = ? WHERE UUID = ?`, UserSession.Following+1, user_id)

	// Mise a jour de la liste des personnes que l'on follow
	db.Exec(`UPDATE users SET followingList = ? WHERE UUID = ?`, UserSession.FollowingList+","+username, user_id)

	// Mise a jour du nombre de followers de la personne que l'on follow
	db.Exec(`UPDATE users SET followers = ? WHERE UUID = ?`, userToFollow.Follower+1, userToFollow.User_id)

	// Mise a jour de la liste des followers de la personne que l'on follow
	db.Exec(`UPDATE users SET followersList = ? WHERE UUID = ?`, userToFollow.FollowerList+","+UserSession.Username, userToFollow.User_id)

	CreateNotification(db, Notification{
		User_id:    userToFollow.User_id,
		User_id2:   UserSession.User_id,
		Posts_id:   "",
		Comment_id: "",
		Reason:     "follow",
		Checked:    false,
	})
}

func UpdateUnfollowing(db *sql.DB, user_id string, username string) { // username etant la personne que l'on va unfollow

	userToUnfollow := GetAccountByUsername(db, username)

	// Mise a jour de notre nombre de following
	db.Exec(`UPDATE users SET following = ? WHERE UUID = ?`, UserSession.Following-1, user_id)

	// Mise a jour de la liste des personnes que l'on unfollow
	db.Exec(`UPDATE users SET followingList = ? WHERE UUID = ?`, strings.Replace(UserSession.FollowingList, ","+username, "", -1), user_id)

	// Mise a jour du nombre de followers de la personne que l'on unfollow
	db.Exec(`UPDATE users SET followers = ? WHERE UUID = ?`, userToUnfollow.Follower-1, userToUnfollow.User_id)

	// Mise a jour de la liste des followers de la personne que l'on unfollow
	db.Exec(`UPDATE users SET followersList = ? WHERE UUID = ?`, strings.Replace(userToUnfollow.FollowerList, ","+UserSession.Username, "", -1), userToUnfollow.User_id)
}

func GetAllMail(db *sql.DB) ([]string, error) {
	query := "SELECT email FROM users"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		mails = append(mails, email)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return mails, nil
}

func GetAllUsers(db *sql.DB) []User {
	query := "SELECT UUID, role, username, email, password, created_at, updated_at, profilePicture, bio, links, categoriesSub, followers, followersList, following, followingList FROM users"
	rows, err := db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.User_id, &user.Role, &user.Username, &user.Email, &user.Password, &user.CreationDate, &user.UpdateDate, &user.Pfp, &user.Bio, &user.Links, &user.CategorieSub, &user.Follower, &user.FollowerList, &user.Following, &user.FollowingList); err != nil {
			return nil
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil
	}

	return users
}

func GetAllDatas(r *http.Request) DataStruct {

	data, _ := getSessionData(r)
	Color := data.User.ColorMode

	return DataStruct{
		User:            UserSession,
		UserTarget:      User{},
		AllUsers:        GetAllUsers(Db),
		RecommendedUser: RecommendedUser{},
		Post:            Post{},
		AllPosts:        GetAllPosts(Db),
		Comment:         Comment{},
		// AllComments:      GetAllComments(),
		Notification: Notification{},
		// AllNotifications: GetAllNotifications(),
		Categorie:     Category{},
		AllCategories: GetAllCategories(Db),
		ColorMode:     Color,
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
