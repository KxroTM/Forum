package forum

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type Post struct {
	Posts_id  string
	User_id   string
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

func UpdatePostDb(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		fmt.Printf("erreur lors de la récupération des utilisateurs depuis la base de données: %v", err)
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.Posts_id, &post.User_id, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker)
		if err != nil {
			fmt.Printf("erreur lors de la lecture des données post depuis la base de données: %v", err)
			continue
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("erreur lors de la récupération des utilisateurs depuis la base de données: %v", err)
	}

	AllPosts = posts
}

func CreatePost(db *sql.DB, user_id string, categorie string, title string, text string) {
	currentTime := time.Now()
	date := currentTime.Format("02-01-2006 15:04")

	uuid, err := uuid.NewV4()
	if err != nil {
		fmt.Println("Erreur lors de la génération de l'UUID :", err)
		return
	}

	PostSession = Post{
		Posts_id:  uuid.String(),
		User_id:   user_id,
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

	db.Exec(`INSERT INTO posts (posts_id, UUID, categorie, title, text, like, liker, dislike, retweet, retweeter, date, report, disliker) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
					`, PostSession.Posts_id, PostSession.User_id, PostSession.Categorie, PostSession.Title, PostSession.Text, PostSession.Like, PostSession.Liker, PostSession.Dislike, PostSession.Retweet, PostSession.Retweeter, PostSession.Date, PostSession.Report, PostSession.Disliker)

	AllPosts = append(AllPosts, PostSession)
	UpdatePostDb(db)
	fmt.Println(AllPosts)
}

func UpdatePost(db *sql.DB, post_id string, categorie string, title string, text string) {
	PostSession = GetPost(db, post_id)
	PostSession.Categorie = categorie
	PostSession.Title = title
	PostSession.Text = text
	db.Exec(`UPDATE posts SET categorie = ?, title = ?, text = ? WHERE posts_id = ?`, categorie, title, text, post_id)
	UpdatePostDb(db)
}

func GetPost(db *sql.DB, post_id string) Post {
	rows, err := db.Query("SELECT * FROM posts WHERE posts_id = ?", post_id)
	if err != nil {
		fmt.Printf("erreur lors de la récupération des utilisateurs depuis la base de données: %v", err)
	}
	defer rows.Close()

	var post Post

	for rows.Next() {
		if err := rows.Scan(&post.Posts_id, &post.User_id, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker); err != nil {
			fmt.Printf("erreur lors de la lecture des données post depuis la base de données: %v", err)
		}
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("erreur lors de la récupération des utilisateurs depuis la base de données: %v", err)
	}

	return post
}

func DeletePost(db *sql.DB, post_id string) {
	db.Exec(`DELETE FROM posts WHERE posts_id = ?`, post_id)
	UpdatePostDb(db)
}

func GetAllPostsByUser(db *sql.DB, user_id string) []Post {
	rows, err := db.Query("SELECT * FROM posts WHERE UUID = ?", user_id)
	if err != nil {
		fmt.Printf("erreur lors de la récupération des posts depuis la base de données: %v", err)
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.Posts_id, &post.User_id, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker)
		if err != nil {
			fmt.Printf("erreur lors de la lecture des données post depuis la base de données: %v", err)
			continue // Passer à l'itération suivante si une erreur se produit
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("erreur lors de la récupération des posts depuis la base de données: %v", err)
	}

	return posts
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

func GetAllPostsByCategorie(db *sql.DB, categorie string) []Post {
	rows, err := db.Query("SELECT * FROM posts WHERE categorie = ?", categorie)
	if err != nil {
		fmt.Printf("erreur lors de la récupération des utilisateurs depuis la base de données: %v", err)
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Posts_id, &post.User_id, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker); err != nil {
			fmt.Printf("erreur lors de la lecture des données post depuis la base de données: %v", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("erreur lors de la récupération des utilisateurs depuis la base de données: %v", err)
	}

	return posts
}

func GetAllPostsByLikeCount(db *sql.DB) []Post {
	rows, err := db.Query("SELECT * FROM posts ORDER BY like DESC")
	if err != nil {
		fmt.Printf("erreur lors de la récupération des utilisateurs depuis la base de données: %v", err)
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Posts_id, &post.User_id, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker); err != nil {
			fmt.Printf("erreur lors de la lecture des données post depuis la base de données: %v", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("erreur lors de la récupération des utilisateurs depuis la base de données: %v", err)
	}

	return posts
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

func LikePost(db *sql.DB, post_id string, username string) {
	PostSession = GetPost(db, post_id)
	PostSession.Liker = PostSession.Liker + "," + username

	if strings.Contains(PostSession.Disliker, username) {
		PostSession.Disliker = strings.Replace(PostSession.Disliker, ","+username, "", -1)

		db.Exec(`UPDATE posts SET like = ?, liker = ?, dislike = ?, disliker = ? WHERE posts_id = ?`, PostSession.Like+1, PostSession.Liker, PostSession.Dislike-1, post_id, PostSession.Disliker)

		UpdatePostDb(db)
	} else {
		db.Exec(`UPDATE posts SET like = ?, liker = ? WHERE posts_id = ?`, PostSession.Like+1, post_id, PostSession.Liker)
		UpdatePostDb(db)
	}
}

func DislikePost(db *sql.DB, post_id string, username string) {
	PostSession = GetPost(db, post_id)
	PostSession.Disliker = PostSession.Disliker + "," + username
	if strings.Contains(PostSession.Liker, username) {
		PostSession.Liker = strings.Replace(PostSession.Liker, ","+username, "", -1)

		db.Exec(`UPDATE posts SET like = ?, liker = ?, dislike = ?, disliker = ? WHERE posts_id = ?`, PostSession.Like-1, PostSession.Liker, PostSession.Dislike+1, post_id, PostSession.Disliker)

		UpdatePostDb(db)
	} else {
		db.Exec(`UPDATE posts SET dislike = ?, disliker = ? WHERE posts_id = ?`, PostSession.Dislike+1, post_id, PostSession.Disliker)
		UpdatePostDb(db)
	}
}

func RetweetPost(db *sql.DB, post_id string, username string) {
	PostSession = GetPost(db, post_id)
	PostSession.Retweeter = PostSession.Retweeter + "," + username
	db.Exec(`UPDATE posts SET retweet = ?, retweeter = ? WHERE posts_id = ?`, PostSession.Retweet+1, PostSession.Retweeter, post_id)
	UpdatePostDb(db)
}

func UnRetweetPost(db *sql.DB, post_id string, username string) {
	PostSession = GetPost(db, post_id)
	PostSession.Retweeter = strings.Replace(PostSession.Retweeter, ","+username, "", -1)
	db.Exec(`UPDATE posts SET retweet = ?, retweeter = ? WHERE posts_id = ?`, PostSession.Retweet-1, PostSession.Retweeter, post_id)
	UpdatePostDb(db)
}

func ReportPost(db *sql.DB, post_id string) {
	PostSession = GetPost(db, post_id)
	db.Exec(`UPDATE posts SET report = ? WHERE posts_id = ?`, PostSession.Report+1, post_id)
	UpdatePostDb(db)
}

func GetPostByReport(db *sql.DB) []Post {
	rows, err := db.Query("SELECT * FROM posts WHERE report > 10")
	if err != nil {
		fmt.Printf("erreur lors de la récupération des utilisateurs depuis la base de données: %v", err)
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Posts_id, &post.User_id, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker); err != nil {
			fmt.Printf("erreur lors de la lecture des données post depuis la base de données: %v", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("erreur lors de la récupération des utilisateurs depuis la base de données: %v", err)
	}

	return posts
}
