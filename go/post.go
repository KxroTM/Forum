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
	IsLike    bool
	IsDislike bool
	IsRetweet bool
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
		rows.Scan(&post.Posts_id, &post.User_id, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker, &post.User_pfp, &post.Author, &post.Links, &post.IsLike, &post.IsDislike, &post.IsRetweet)
		post.Links = strings.TrimSpace(post.Links)
		if UserSession.Username != "" {
			if strings.Contains(post.Liker, UserSession.Username) {
				post.IsLike = true
			} else {
				post.IsLike = false
			}
			if strings.Contains(post.Disliker, UserSession.Username) {
				post.IsDislike = true
			} else {
				post.IsDislike = false
			}
			if strings.Contains(post.Retweeter, UserSession.Username) {
				post.IsRetweet = true
			} else {
				post.IsRetweet = false
			}
		}
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
		IsLike:    false,
		IsDislike: false,
		IsRetweet: false,
	}

	db.Exec(`INSERT INTO posts (posts_id, UUID, categorie, title, text, like, liker, dislike, retweet, retweeter, date, report, disliker, user_pfp, author, links, isLike, isDislike, isRetweet) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
					`, PostSession.Posts_id, PostSession.User_id, PostSession.Categorie, PostSession.Title, PostSession.Text, PostSession.Like, PostSession.Liker, PostSession.Dislike, PostSession.Retweet, PostSession.Retweeter, PostSession.Date, PostSession.Report, PostSession.Disliker, PostSession.User_pfp, PostSession.Author, PostSession.Links, PostSession.IsLike, PostSession.IsDislike, PostSession.IsRetweet)
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
		rows.Scan(&post.Posts_id, &post.User_id, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker, &post.User_pfp, &post.Author, &post.Links, &post.IsLike, &post.IsDislike, &post.IsRetweet)
	}
	if err := rows.Err(); err != nil {
		return Post{}, err
	}
	post.Links = strings.TrimSpace(post.Links)
	if UserSession.Username != "" {
		if strings.Contains(post.Liker, UserSession.Username) {
			post.IsLike = true
		} else {
			post.IsLike = false
		}
		if strings.Contains(post.Disliker, UserSession.Username) {
			post.IsDislike = true
		} else {
			post.IsDislike = false
		}
		if strings.Contains(post.Retweeter, UserSession.Username) {
			post.IsRetweet = true
		} else {
			post.IsRetweet = false
		}
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
		rows.Scan(&post.Posts_id, &post.User_id, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker, &post.User_pfp, &post.Author, &post.Links, &post.IsLike, &post.IsDislike, &post.IsRetweet)
		post.Links = strings.TrimSpace(post.Links)
		if UserSession.Username != "" {
			if strings.Contains(post.Liker, UserSession.Username) {
				post.IsLike = true
			} else {
				post.IsLike = false
			}
			if strings.Contains(post.Disliker, UserSession.Username) {
				post.IsDislike = true
			} else {
				post.IsDislike = false
			}
			if strings.Contains(post.Retweeter, UserSession.Username) {
				post.IsRetweet = true
			} else {
				post.IsRetweet = false
			}
		}
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
	rows, err := db.Query("SELECT * FROM posts WHERE categorie LIKE ?", "%"+categorie+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		rows.Scan(&post.Posts_id, &post.User_id, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker, &post.User_pfp, &post.Author, &post.Links, &post.IsLike, &post.IsDislike, &post.IsRetweet)
		post.Links = strings.TrimSpace(post.Links)
		if UserSession.Username != "" {
			if strings.Contains(post.Liker, UserSession.Username) {
				post.IsLike = true
			} else {
				post.IsLike = false
			}
			if strings.Contains(post.Disliker, UserSession.Username) {
				post.IsDislike = true
			} else {
				post.IsDislike = false
			}
			if strings.Contains(post.Retweeter, UserSession.Username) {
				post.IsRetweet = true
			} else {
				post.IsRetweet = false
			}
		}
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
		rows.Scan(&post.Posts_id, &post.User_id, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker, &post.User_pfp, &post.Author, &post.Links, &post.IsLike, &post.IsDislike, &post.IsRetweet)
		post.Links = strings.TrimSpace(post.Links)
		if UserSession.Username != "" {
			if strings.Contains(post.Liker, UserSession.Username) {
				post.IsLike = true
			} else {
				post.IsLike = false
			}
			if strings.Contains(post.Disliker, UserSession.Username) {
				post.IsDislike = true
			} else {
				post.IsDislike = false
			}
			if strings.Contains(post.Retweeter, UserSession.Username) {
				post.IsRetweet = true
			} else {
				post.IsRetweet = false
			}
		}
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
			post.Links = strings.TrimSpace(post.Links)
			if UserSession.Username != "" {
				if strings.Contains(post.Liker, UserSession.Username) {
					post.IsLike = true
				} else {
					post.IsLike = false
				}
				if strings.Contains(post.Disliker, UserSession.Username) {
					post.IsDislike = true
				} else {
					post.IsDislike = false
				}
				if strings.Contains(post.Retweeter, UserSession.Username) {
					post.IsRetweet = true
				} else {
					post.IsRetweet = false
				}
			}
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
			post.Links = strings.TrimSpace(post.Links)
			if UserSession.Username != "" {
				if strings.Contains(post.Liker, UserSession.Username) {
					post.IsLike = true
				} else {
					post.IsLike = false
				}
				if strings.Contains(post.Disliker, UserSession.Username) {
					post.IsDislike = true
				} else {
					post.IsDislike = false
				}
				if strings.Contains(post.Retweeter, UserSession.Username) {
					post.IsRetweet = true
				} else {
					post.IsRetweet = false
				}
			}
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

	if PostSession.Liker == "" {
		PostSession.Liker = username
	} else {
		PostSession.Liker = PostSession.Liker + "," + username
	}

	if strings.Contains(PostSession.Disliker, username) {
		UnDislikePost(db, post_id, username)

	}
	db.Exec(`UPDATE posts SET like = ?, liker = ? WHERE posts_id = ?`, PostSession.Like+1, PostSession.Liker, post_id)

	err = CreateNotification(db, Notification{
		User_id:    PostSession.User_id,
		User_id2:   UserSession.User_id,
		Posts_id:   post_id,
		Comment_id: "",
		Reason:     "likePost",
		Checked:    false,
	})
	if err != nil {
		return err
	}

	return nil
}

func UnLikePost(db *sql.DB, post_id string, username string) error {
	PostSession, err := GetPost(db, post_id)
	if err != nil {
		return err
	}
	PostSession.Liker = strings.Replace(PostSession.Liker, username, "", -1)
	db.Exec(`UPDATE posts SET like = ?, liker = ? WHERE posts_id = ?`, PostSession.Like-1, PostSession.Liker, post_id)
	return nil
}

func DislikePost(db *sql.DB, post_id string, username string) error {
	PostSession, err := GetPost(db, post_id)
	if err != nil {
		return err
	}

	if PostSession.Disliker == "" {
		PostSession.Disliker = username
	} else {
		PostSession.Disliker = PostSession.Disliker + "," + username
	}

	if strings.Contains(PostSession.Liker, username) {
		UnLikePost(db, post_id, username)

	}
	db.Exec(`UPDATE posts SET dislike = ?, disliker = ? WHERE posts_id = ?`, PostSession.Dislike+1, PostSession.Disliker, post_id)

	return nil
}

func UnDislikePost(db *sql.DB, post_id string, username string) error {
	PostSession, err := GetPost(db, post_id)
	if err != nil {
		return err
	}
	PostSession.Disliker = strings.Replace(PostSession.Disliker, username, "", -1)
	db.Exec(`UPDATE posts SET dislike = ?, disliker = ? WHERE posts_id = ?`, PostSession.Dislike-1, PostSession.Disliker, post_id)
	return nil
}

func RetweetPost(db *sql.DB, post_id string, username string) error {
	PostSession, err := GetPost(db, post_id)
	if err != nil {
		return err
	}
	if PostSession.Retweeter == "" {
		PostSession.Retweeter = username
	} else {
		PostSession.Retweeter = PostSession.Retweeter + "," + username
	}

	db.Exec(`UPDATE posts SET retweet = ?, retweeter = ? WHERE posts_id = ?`, PostSession.Retweet+1, PostSession.Retweeter, post_id)

	err = CreateNotification(db, Notification{
		User_id:    PostSession.User_id,
		User_id2:   UserSession.User_id,
		Posts_id:   post_id,
		Comment_id: "",
		Reason:     "repost",
		Checked:    false,
	})
	if err != nil {
		return err
	}

	return nil
}

func UnRetweetPost(db *sql.DB, post_id string, username string) error {
	PostSession, err := GetPost(db, post_id)
	if err != nil {
		return err
	}
	PostSession.Retweeter = strings.Replace(PostSession.Retweeter, username, "", -1)
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
		rows.Scan(&post.Posts_id, &post.User_id, &post.Categorie, &post.Title, &post.Text, &post.Like, &post.Liker, &post.Dislike, &post.Retweet, &post.Retweeter, &post.Date, &post.Report, &post.Disliker, &post.User_pfp, &post.Author, &post.Links, &post.IsLike, &post.IsDislike, &post.IsRetweet)
		post.Links = strings.TrimSpace(post.Links)
		if UserSession.Username != "" {
			if strings.Contains(post.Liker, UserSession.Username) {
				post.IsLike = true
			} else {
				post.IsLike = false
			}
			if strings.Contains(post.Disliker, UserSession.Username) {
				post.IsDislike = true
			} else {
				post.IsDislike = false
			}
			if strings.Contains(post.Retweeter, UserSession.Username) {
				post.IsRetweet = true
			} else {
				post.IsRetweet = false
			}
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostByFollowing(db *sql.DB, user_id string) []Post {
	var posts []Post
	var AllPosts = GetAllPosts(db)
	user := GetAccountById(db, user_id)

	for _, post := range AllPosts {
		if strings.Contains(user.FollowingList, post.Author) {
			post.Links = strings.TrimSpace(post.Links)
			if UserSession.Username != "" {
				if strings.Contains(post.Liker, UserSession.Username) {
					post.IsLike = true
				} else {
					post.IsLike = false
				}
				if strings.Contains(post.Disliker, UserSession.Username) {
					post.IsDislike = true
				} else {
					post.IsDislike = false
				}
				if strings.Contains(post.Retweeter, UserSession.Username) {
					post.IsRetweet = true
				} else {
					post.IsRetweet = false
				}
			}
			posts = append(posts, post)
		}
	}
	return posts
}

func GetPostBySearch(db *sql.DB, search string, from []Post) []Post {
	var posts []Post
	var AllPosts = from

	for _, post := range AllPosts {
		if strings.Contains(strings.ToLower(post.Title), strings.ToLower(search)) || strings.Contains(strings.ToLower(post.Text), strings.ToLower(search)) || strings.Contains(strings.ToLower(post.Author), strings.ToLower(search)) {
			post.Links = strings.TrimSpace(post.Links)
			if UserSession.Username != "" {
				if strings.Contains(post.Liker, UserSession.Username) {
					post.IsLike = true
				} else {
					post.IsLike = false
				}
				if strings.Contains(post.Disliker, UserSession.Username) {
					post.IsDislike = true
				} else {
					post.IsDislike = false
				}
				if strings.Contains(post.Retweeter, UserSession.Username) {
					post.IsRetweet = true
				} else {
					post.IsRetweet = false
				}
			}
			posts = append(posts, post)
		}
	}
	return posts
}
