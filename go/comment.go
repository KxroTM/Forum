package forum

import (
	"database/sql"
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
