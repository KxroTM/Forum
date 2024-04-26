package forum

var Post struct {
	posts_id  int
	user_id   int
	categorie string
	title     string
	text      string
	like      int
	liker     []int
	dislike   int
	retweet   int
	retweeter []int
	date      string
	report    int
}
