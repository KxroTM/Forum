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
	Author    string
	Links     string
}

var PostSession Post

func GetAllPosts(db *sql.DB) []Post {
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		rows.Scan(&post.Posts_id, &post.User_id, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker, &post.User_pfp, &post.Author, &post.Links)
		post.Links = strings.TrimSpace(post.Links)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return posts
}

func CreatePost(db *sql.DB, user_id string, categorie string, title string, text string, link string) error {
	currentTime := time.Now()
	date := currentTime.Format("02-01-2006 15:04")

	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	user := GetAccountById(db, user_id)

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
		Author:    user.Username,
		Links:     link,
	}

	db.Exec(`INSERT INTO posts (posts_id, UUID, user_pfp, categorie, title, text, like, liker, dislike, retweet, retweeter, date, report, disliker, author, links) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
					`, PostSession.Posts_id, PostSession.User_id, PostSession.User_pfp, PostSession.Categorie, PostSession.Title, PostSession.Text, PostSession.Like, PostSession.Liker, PostSession.Dislike, PostSession.Retweet, PostSession.Retweeter, PostSession.Date, PostSession.Report, PostSession.Disliker, PostSession.Author, PostSession.Links)
	return nil
}

func UpdatePost(db *sql.DB, post_id string, categorie string, title string, text string, links string) error {
	PostSession, err := GetPost(db, post_id)
	if err != nil {
		return err
	}
	PostSession.Categorie = categorie
	PostSession.Title = title
	PostSession.Text = text
	PostSession.Links = links
	db.Exec(`UPDATE posts SET categorie = ?, title = ?, text = ?, links = ? WHERE posts_id = ?`, categorie, title, text, links, post_id)
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
		rows.Scan(&post.Posts_id, &post.User_id, &post.User_pfp, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker, &post.Author, &post.Links)
	}
	if err := rows.Err(); err != nil {
		return Post{}, err
	}

	return post, nil
}

func DeletePost(db *sql.DB, post_id string) {
	db.Exec(`DELETE FROM posts WHERE posts_id = ?`, post_id)
}

func DeleteAllPosts(db *sql.DB) {
	db.Exec(`DELETE FROM posts`)
}

func DeleteAllPostsByUser(db *sql.DB, user_id string) {
	db.Exec(`DELETE FROM posts WHERE UUID = ?`, user_id)
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
		rows.Scan(&post.Posts_id, &post.User_id, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker, &post.User_pfp, &post.Author, &post.Links)
		post.Links = strings.TrimSpace(post.Links)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetAllPostsByDate(db *sql.DB) []Post {
	var allPostsTemp []Post
	var AllPosts = GetAllPosts(db)

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
		rows.Scan(&post.Posts_id, &post.User_id, &post.User_pfp, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker, &post.Author, &post.Links)
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
		rows.Scan(&post.Posts_id, &post.User_id, &post.User_pfp, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker, &post.Author, &post.Links)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetAllPostsByRetweet(db *sql.DB, username string) []Post {
	var posts []Post
	var AllPosts = GetAllPosts(db)

	for _, post := range AllPosts {
		if strings.Contains(post.Retweeter, username) {
			posts = append(posts, post)
		}
	}

	return posts
}

func GetAllPostsByLike(db *sql.DB, username string) []Post {
	var posts []Post
	var AllPosts = GetAllPosts(db)

	for _, post := range AllPosts {
		if strings.Contains(post.Liker, username) {
			posts = append(posts, post)
		}
	}

	return posts
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

	} else {
		db.Exec(`UPDATE posts SET like = ?, liker = ? WHERE posts_id = ?`, PostSession.Like+1, post_id, PostSession.Liker)
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

	} else {
		db.Exec(`UPDATE posts SET dislike = ?, disliker = ? WHERE posts_id = ?`, PostSession.Dislike+1, post_id, PostSession.Disliker)

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

	return nil
}

func UnRetweetPost(db *sql.DB, post_id string, username string) error {
	PostSession, err := GetPost(db, post_id)
	if err != nil {
		return err
	}
	PostSession.Retweeter = strings.Replace(PostSession.Retweeter, ","+username, "", -1)
	db.Exec(`UPDATE posts SET retweet = ?, retweeter = ? WHERE posts_id = ?`, PostSession.Retweet-1, PostSession.Retweeter, post_id)
	return nil
}

func ReportPost(db *sql.DB, post_id string) error {
	PostSession, err := GetPost(db, post_id)
	if err != nil {
		return err
	}
	db.Exec(`UPDATE posts SET report = ? WHERE posts_id = ?`, PostSession.Report+1, post_id)
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
		rows.Scan(&post.Posts_id, &post.User_id, &post.User_pfp, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker, &post.Author)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
