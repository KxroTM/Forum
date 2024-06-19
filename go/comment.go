package forum

import (
	"database/sql"
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
}

func GetAllComment(db *sql.DB) []Comment {
	rows, err := db.Query("SELECT * FROM comments")
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
