package forum

import (
	"database/sql"
	"fmt"
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

func UpdateCommentDb(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM comments")
	if err != nil {
		fmt.Printf("erreur lors de la récupération des utilisateurs depuis la base de données: %v", err)
	}
	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.Comment_id, &comment.Posts_id, &comment.User_id, &comment.Text, &comment.Date, &comment.Like, &comment.Dislike, &comment.Report, &comment.Liker, &comment.Disliker)
		if err != nil {
			fmt.Printf("erreur lors de la lecture des données post depuis la base de données: %v", err)
			continue
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("erreur lors de la récupération des utilisateurs depuis la base de données: %v", err)
	}

	AllComments = comments
}
