package forum

var Comment struct {
	comment_id int
	posts_id   int
	user_id    int
	text       string
	date       string
	like       int
	dislike    int
	report     int
}
