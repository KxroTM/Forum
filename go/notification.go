package forum

import "fmt"

type Notification struct {
	Notification_id string
	Comment_id      string
	Posts_id        string
	User_id         string
	User_id2        string
	Date            string
	Checked         bool
}

func NewFollower(Notification Notification) { //la fonction qui permet de notifier un utilisateur qu'il a un nouveau follower
	if UserSession.User_id != Notification.User_id { //si l'utilisateur n'est pas celui qui a un nouveau follower
		fmt.Println("You have a new follower:", Notification.User_id) //on affiche le message suivant à l'utilisateur qui a un nouveau follower : "You have a new follower:" suivi du nom de l'utilisateur qui le suit
	}
}
