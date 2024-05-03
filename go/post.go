package forum

import (
	"database/sql"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type Post struct {
	Posts_id  string
	User_id   string
	User_pfp  string
	Categorie string
	Title     string
	Text      string
	Like      int
	Liker     string
	Dislike   int
	Disliker  string
	Retweet   int
	Retweeter string
	Date      string
	Report    int
}

var PostSession Post
var AllPosts []Post

func UpdatePostDb(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		return err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		rows.Scan(&post.Posts_id, &post.User_id, &post.User_pfp, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker)
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	AllPosts = posts
	return nil
}

func CreatePost(db *sql.DB, user_id string, categorie string, title string, text string) error {
	currentTime := time.Now()
	date := currentTime.Format("02-01-2006 15:04")

	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	user := GetAccountById(user_id)

	PostSession = Post{
		Posts_id:  uuid.String(),
		User_id:   user_id,
		User_pfp:  user.Pfp,
		Categorie: categorie,
		Title:     title,
		Text:      text,
		Like:      0,
		Liker:     "",
		Dislike:   0,
		Retweet:   0,
		Retweeter: "",
		Date:      date,
		Report:    0,
		Disliker:  "",
	}

	db.Exec(`INSERT INTO posts (posts_id, UUID, user_pfp, categorie, title, text, like, liker, dislike, retweet, retweeter, date, report, disliker) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
					`, PostSession.Posts_id, PostSession.User_id, PostSession.User_pfp, PostSession.Categorie, PostSession.Title, PostSession.Text, PostSession.Like, PostSession.Liker, PostSession.Dislike, PostSession.Retweet, PostSession.Retweeter, PostSession.Date, PostSession.Report, PostSession.Disliker)

	AllPosts = append(AllPosts, PostSession)
	UpdatePostDb(db)
	return nil
}

func UpdatePost(db *sql.DB, post_id string, categorie string, title string, text string) error {
	PostSession, err := GetPost(db, post_id)
	if err != nil {
		return err
	}
	PostSession.Categorie = categorie
	PostSession.Title = title
	PostSession.Text = text
	db.Exec(`UPDATE posts SET categorie = ?, title = ?, text = ? WHERE posts_id = ?`, categorie, title, text, post_id)
	UpdatePostDb(db)
	return nil
}

func GetPost(db *sql.DB, post_id string) (Post, error) {
	rows, err := db.Query("SELECT * FROM posts WHERE posts_id = ?", post_id)
	if err != nil {
		return Post{}, err
	}
	defer rows.Close()

	var post Post

	for rows.Next() {
		rows.Scan(&post.Posts_id, &post.User_id, &post.User_pfp, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker)
	}
	if err := rows.Err(); err != nil {
		return Post{}, err
	}

	return post, nil
}

func DeletePost(db *sql.DB, post_id string) {
	db.Exec(`DELETE FROM posts WHERE posts_id = ?`, post_id)
	UpdatePostDb(db)
}

func GetAllPostsByUser(db *sql.DB, user_id string) ([]Post, error) {
	rows, err := db.Query("SELECT * FROM posts WHERE UUID = ?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		rows.Scan(&post.Posts_id, &post.User_id, &post.User_pfp, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker)
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetAllPostsByDate() []Post {
	var allPostsTemp []Post

	allPostsTemp = append(allPostsTemp, AllPosts...)

	sort.Slice(allPostsTemp, func(i, j int) bool {
		date1, _ := time.Parse("02-01-2006 15:04", allPostsTemp[i].Date)
		date2, _ := time.Parse("02-01-2006 15:04", allPostsTemp[j].Date)
		return date1.After(date2)
	})

	return allPostsTemp
}

func GetAllPostsByCategorie(db *sql.DB, categorie string) ([]Post, error) {
	rows, err := db.Query("SELECT * FROM posts WHERE categorie = ?", categorie)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		rows.Scan(&post.Posts_id, &post.User_id, &post.User_pfp, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetAllPostsByLikeCount(db *sql.DB) ([]Post, error) {
	rows, err := db.Query("SELECT * FROM posts ORDER BY like DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		rows.Scan(&post.Posts_id, &post.User_id, &post.User_pfp, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetAllPostsByRetweet(username string) []Post {
	var posts []Post

	for _, post := range AllPosts {
		if strings.Contains(post.Retweeter, username) {
			posts = append(posts, post)
		}
	}

	return posts
}

func GetAllPostsByLike(username string) []Post {
	var posts []Post

	for _, post := range AllPosts {
		if strings.Contains(post.Liker, username) {
			posts = append(posts, post)
		}
	}

	return posts
}

func GetAllPosts() []Post {
	return AllPosts
}

func LikePost(db *sql.DB, post_id string, username string) error {
	PostSession, err := GetPost(db, post_id)
	if err != nil {
		return err
	}
	PostSession.Liker = PostSession.Liker + "," + username

	if strings.Contains(PostSession.Disliker, username) {
		PostSession.Disliker = strings.Replace(PostSession.Disliker, ","+username, "", -1)

		db.Exec(`UPDATE posts SET like = ?, liker = ?, dislike = ?, disliker = ? WHERE posts_id = ?`, PostSession.Like+1, PostSession.Liker, PostSession.Dislike-1, post_id, PostSession.Disliker)

		UpdatePostDb(db)
	} else {
		db.Exec(`UPDATE posts SET like = ?, liker = ? WHERE posts_id = ?`, PostSession.Like+1, post_id, PostSession.Liker)
		UpdatePostDb(db)
	}
	return nil
}

func DislikePost(db *sql.DB, post_id string, username string) error {
	PostSession, err := GetPost(db, post_id)
	if err != nil {
		return err
	}
	PostSession.Disliker = PostSession.Disliker + "," + username
	if strings.Contains(PostSession.Liker, username) {
		PostSession.Liker = strings.Replace(PostSession.Liker, ","+username, "", -1)

		db.Exec(`UPDATE posts SET like = ?, liker = ?, dislike = ?, disliker = ? WHERE posts_id = ?`, PostSession.Like-1, PostSession.Liker, PostSession.Dislike+1, post_id, PostSession.Disliker)

		UpdatePostDb(db)
	} else {
		db.Exec(`UPDATE posts SET dislike = ?, disliker = ? WHERE posts_id = ?`, PostSession.Dislike+1, post_id, PostSession.Disliker)
		UpdatePostDb(db)
	}
	return nil
}

func RetweetPost(db *sql.DB, post_id string, username string) error {
	PostSession, err := GetPost(db, post_id)
	if err != nil {
		return err
	}
	PostSession.Retweeter = PostSession.Retweeter + "," + username
	db.Exec(`UPDATE posts SET retweet = ?, retweeter = ? WHERE posts_id = ?`, PostSession.Retweet+1, PostSession.Retweeter, post_id)
	UpdatePostDb(db)
	return nil
}

func UnRetweetPost(db *sql.DB, post_id string, username string) error {
	PostSession, err := GetPost(db, post_id)
	if err != nil {
		return err
	}
	PostSession.Retweeter = strings.Replace(PostSession.Retweeter, ","+username, "", -1)
	db.Exec(`UPDATE posts SET retweet = ?, retweeter = ? WHERE posts_id = ?`, PostSession.Retweet-1, PostSession.Retweeter, post_id)
	UpdatePostDb(db)
	return nil
}

func ReportPost(db *sql.DB, post_id string) error {
	PostSession, err := GetPost(db, post_id)
	if err != nil {
		return err
	}
	db.Exec(`UPDATE posts SET report = ? WHERE posts_id = ?`, PostSession.Report+1, post_id)
	UpdatePostDb(db)
	return nil
}

func GetPostByReport(db *sql.DB) ([]Post, error) {
	rows, err := db.Query("SELECT * FROM posts WHERE report > 10")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		rows.Scan(&post.Posts_id, &post.User_id, &post.User_pfp, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
