package forum

import (
	"database/sql"
)

type Comment struct {
	Comment_id string
	Posts_id   string
	User_id    string
	User_pfp   string
	Text       string
	Date       string
	Like       int
	Dislike    int
	Report     int
	Liker      string
	Disliker   string
}

var CommentSession Comment
var AllComments []Comment

func UpdateCommentDb(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM comments")
	if err != nil {
		return err
	}
	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var comment Comment
		rows.Scan(&comment.Comment_id, &comment.Posts_id, &comment.User_id, &comment.Text, &comment.Date, &comment.Like, &comment.Dislike, &comment.Report, &comment.Liker, &comment.Disliker)
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	AllComments = comments
	return nil
}
