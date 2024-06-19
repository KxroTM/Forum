package forum

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type Comment struct {
	Comment_id string
	Posts_id   string
	User_id    string
	Text       string
	Date       string
	Like       int
	Dislike    int
	Report     int
	Liker      string
	Disliker   string
	User_pfp   string
	Author     string
	PfpChanged bool
	IsLike     bool
	IsDislike  bool
}

func GetAllCommentByUser(db *sql.DB, user_id string) []Comment {
	rows, err := db.Query("SELECT * FROM comments WHERE UUID = ?", user_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.Comment_id, &comment.Posts_id, &comment.User_id, &comment.Text, &comment.Date, &comment.Like, &comment.Dislike, &comment.Report, &comment.Liker, &comment.Disliker, &comment.User_pfp)
		if err != nil {
			panic(err)
		}
		user := GetAccountById(db, comment.User_id)
		comment.Author = user.Username
		if !strings.Contains(user.Pfp, "../../style/media/default_avatar/") {
			comment.PfpChanged = true
		}
		if UserSession.Username != "" {
			if strings.Contains(comment.Liker, UserSession.Username) {
				comment.IsLike = true
			} else {
				comment.IsLike = false
			}
			if strings.Contains(comment.Disliker, UserSession.Username) {
				comment.IsDislike = true
			} else {
				comment.IsDislike = false
			}
		}
		comments = append(comments, comment)
	}
	return comments
}

func GetComment(db *sql.DB, comment_id string) Comment {
	row := db.QueryRow("SELECT * FROM comments WHERE comment_id = ?", comment_id)
	var comment Comment
	err := row.Scan(&comment.Comment_id, &comment.Posts_id, &comment.User_id, &comment.Text, &comment.Date, &comment.Like, &comment.Dislike, &comment.Report, &comment.Liker, &comment.Disliker, &comment.User_pfp)
	if err != nil {
		panic(err)
	}
	user := GetAccountById(db, comment.User_id)
	comment.Author = user.Username
	if !strings.Contains(user.Pfp, "../../style/media/default_avatar/") {
		comment.PfpChanged = true
	}
	if UserSession.Username != "" {
		if strings.Contains(comment.Liker, UserSession.Username) {
			comment.IsLike = true
		} else {
			comment.IsLike = false
		}
		if strings.Contains(comment.Disliker, UserSession.Username) {
			comment.IsDislike = true
		} else {
			comment.IsDislike = false
		}
	}
	return comment
}

func GetCommentByPostId(db *sql.DB, posts_id string) []Comment {
	rows, err := db.Query("SELECT * FROM comments WHERE post_id = ?", posts_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.Comment_id, &comment.Posts_id, &comment.User_id, &comment.Text, &comment.Date, &comment.Like, &comment.Dislike, &comment.Report, &comment.Liker, &comment.Disliker, &comment.User_pfp)
		if err != nil {
			panic(err)
		}

		user := GetAccountById(db, comment.User_id)
		comment.Author = user.Username
		if !strings.Contains(user.Pfp, "../../style/media/default_avatar/") {
			comment.PfpChanged = true
		}
		if UserSession.Username != "" {
			if strings.Contains(comment.Liker, UserSession.Username) {
				comment.IsLike = true
			} else {
				comment.IsLike = false
			}
			if strings.Contains(comment.Disliker, UserSession.Username) {
				comment.IsDislike = true
			} else {
				comment.IsDislike = false
			}
		}
		comments = append(comments, comment)
	}
	return comments
}

func CreateCommentaire(db *sql.DB, commentaire, post_id, user_id string) {
	var Commentaire Comment
	currentTime := time.Now()
	date := currentTime.Format("02-01-2006 15:04")

	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	user := GetAccountById(db, user_id)

	Commentaire = Comment{
		Comment_id: uuid.String(),
		Posts_id:   post_id,
		User_id:    user_id,
		Text:       commentaire,
		Date:       date,
		Like:       0,
		Dislike:    0,
		Report:     0,
		Liker:      "",
		Disliker:   "",
		User_pfp:   user.Pfp,
	}

	_, err = db.Exec("INSERT INTO comments (comment_id, post_id, UUID, text, date, like, dislike, report, liker, disliker, user_pfp) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", Commentaire.Comment_id, Commentaire.Posts_id, Commentaire.User_id, Commentaire.Text, Commentaire.Date, Commentaire.Like, Commentaire.Dislike, Commentaire.Report, Commentaire.Liker, Commentaire.Disliker, Commentaire.User_pfp)
	if err != nil {
		panic(err)
	}
}

func LikeComment(db *sql.DB, comment_id string, username string) error {
	CommentSession := GetComment(db, comment_id)

	if CommentSession.Liker == "" {
		CommentSession.Liker = username
	} else {
		CommentSession.Liker = CommentSession.Liker + "," + username
	}

	if strings.Contains(CommentSession.Disliker, username) {
		UnDislikeComment(db, comment_id, username)

	}
	db.Exec(`UPDATE comments SET like = ?, liker = ? WHERE comment_id = ?`, CommentSession.Like+1, CommentSession.Liker, comment_id)

	err := CreateNotification(db, Notification{
		User_id:    CommentSession.User_id,
		User_id2:   UserSession.User_id,
		Posts_id:   "",
		Comment_id: comment_id,
		Reason:     "likeComment",
		Checked:    false,
	})
	if err != nil {
		return err
	}

	return nil
}

func UnLikeComment(db *sql.DB, comment_id string, username string) error {
	CommentSession := GetComment(db, comment_id)
	CommentSession.Liker = strings.Replace(CommentSession.Liker, username, "", -1)
	db.Exec(`UPDATE comments SET like = ?, liker = ? WHERE comment_id = ?`, CommentSession.Like-1, CommentSession.Liker, comment_id)
	return nil
}

func DislikeComment(db *sql.DB, comment_id string, username string) error {
	CommentSession := GetComment(db, comment_id)

	if CommentSession.Disliker == "" {
		CommentSession.Disliker = username
	} else {
		CommentSession.Disliker = CommentSession.Disliker + "," + username
	}

	if strings.Contains(CommentSession.Liker, username) {
		UnLikeComment(db, comment_id, username)

	}
	db.Exec(`UPDATE comments SET dislike = ?, disliker = ? WHERE comment_id = ?`, CommentSession.Dislike+1, CommentSession.Disliker, comment_id)

	return nil
}

func UnDislikeComment(db *sql.DB, comment_id string, username string) error {
	CommentSession := GetComment(db, comment_id)

	CommentSession.Disliker = strings.Replace(CommentSession.Disliker, username, "", -1)
	db.Exec(`UPDATE comments SET dislike = ?, disliker = ? WHERE comment_id = ?`, CommentSession.Dislike-1, CommentSession.Disliker, comment_id)
	return nil
}

func GetAllPostByComments(db *sql.DB, comments []Comment) []Post {
	var posts []Post
	allPosts := GetAllPosts(Db)
	for _, comment := range comments {
		for _, post := range allPosts {
			if comment.Posts_id == post.Posts_id {
				if !contains(posts, post) {
					posts = append(posts, post)
				}
			}
		}
	}
	return posts
}
