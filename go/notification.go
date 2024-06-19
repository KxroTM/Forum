package forum

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type Notification struct {
	Notification_id string
	Comment_id      string
	Posts_id        string
	User_id         string
	User_id2        string
	Date            string
	Checked         bool
	Reason          string
	User_pfp        string
	PfpChanged      bool
}

func CreateNotification(Db *sql.DB, Notification Notification) error {

	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	Notification.Notification_id = uuid.String()

	currentTime := time.Now()
	date := currentTime.Format("02-01-2006 15:04")

	_, err = Db.Exec("INSERT INTO notifications (notification_id, comment_id, posts_id, user_id, user_id2, date, checked, reason) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", Notification.Notification_id, Notification.Comment_id, Notification.Posts_id, Notification.User_id, Notification.User_id2, date, false, Notification.Reason)

	if err != nil {
		return err
	}

	return err
}

func GetNotifications(Db *sql.DB, User_id string) []Notification {
	rows, err := Db.Query("SELECT * FROM notifications WHERE user_id = ?", User_id)
	if err != nil {
		fmt.Println("Error getting notifications:", err)
	}
	defer rows.Close()
	var Notifications []Notification
	for rows.Next() {
		var Notification Notification
		err := rows.Scan(&Notification.Notification_id, &Notification.Comment_id, &Notification.Posts_id, &Notification.User_id, &Notification.User_id2, &Notification.Date, &Notification.Checked, &Notification.Reason)
		if err != nil {
			fmt.Println("Error getting notifications:", err)
		}
		if GetAccountById(Db, Notification.User_id2).PfpChanged {
			Notification.PfpChanged = true
		}

		Notification.User_pfp = GetAccountById(Db, Notification.User_id2).Pfp
		Notification.User_id2 = GetAccountById(Db, Notification.User_id2).Username

		Notifications = append(Notifications, Notification)
	}
	return Notifications
}

func GetNotification(Db *sql.DB, Notification_id string) Notification {
	row := Db.QueryRow("SELECT * FROM notifications WHERE notification_id = ?", Notification_id)
	var Notification Notification
	err := row.Scan(&Notification.Notification_id, &Notification.Comment_id, &Notification.Posts_id,
		&Notification.User_id, &Notification.User_id2, &Notification.Date, &Notification.Checked, &Notification.Reason)
	if err != nil {
		fmt.Println("Error getting notification:", err)
	}
	if GetAccountById(Db, Notification.User_id2).PfpChanged {
		Notification.PfpChanged = true
	}
	return Notification
}

func GetUnreadNotifications(Db *sql.DB, User_id string) []Notification {
	rows, err := Db.Query("SELECT * FROM notifications WHERE user_id = ? AND checked = FALSE", User_id)
	if err != nil {
		fmt.Println("Error getting notifications:", err)
	}
	defer rows.Close()
	var Notifications []Notification
	for rows.Next() {
		var Notification Notification
		err := rows.Scan(&Notification.Notification_id, &Notification.Comment_id, &Notification.Posts_id,
			&Notification.User_id, &Notification.User_id2, &Notification.Date, &Notification.Checked, &Notification.Reason)
		if err != nil {
			fmt.Println("Error getting notifications:", err)
		}
		if GetAccountById(Db, Notification.User_id2).PfpChanged {
			Notification.PfpChanged = true
		}
		Notifications = append(Notifications, Notification)
	}
	return Notifications
}

func ReadNotification(Db *sql.DB, Notification_id string) {
	_, err := Db.Exec("UPDATE notifications SET checked = TRUE WHERE notification_id = ?", Notification_id)
	if err != nil {
		fmt.Println("Error checking notification:", err)
	}
}

func DeleteNotification(Db *sql.DB, Notification_id string) {
	_, err := Db.Exec("DELETE FROM notifications WHERE notification_id = ?", Notification_id)
	if err != nil {
		fmt.Println("Error deleting notification:", err)
	}
}

func DeleteNotifications(Db *sql.DB, User_id string) {
	_, err := Db.Exec("DELETE FROM notifications WHERE user_id = ?", User_id)
	if err != nil {
		fmt.Println("Error deleting notifications:", err)
	}
}

func DeleteNotificationsByPost(Db *sql.DB, Posts_id string) {
	_, err := Db.Exec("DELETE FROM notifications WHERE posts_id = ?", Posts_id)
	if err != nil {
		fmt.Println("Error deleting notifications:", err)
	}
}

func DeleteNotificationsByComment(Db *sql.DB, Comment_id string) {
	_, err := Db.Exec("DELETE FROM notifications WHERE comment_id = ?", Comment_id)
	if err != nil {
		fmt.Println("Error deleting notifications:", err)
	}
}

func GetCountNotifications(Db *sql.DB, User_id string) int {
	row := Db.QueryRow("SELECT COUNT(*) FROM notifications WHERE user_id = ? AND checked = FALSE", User_id)
	var count int
	err := row.Scan(&count)
	if err != nil {
		fmt.Println("Error getting count notifications:", err)
	}
	return count
}

func GetNotifBySearch(Db *sql.DB, User_id string, Search string) []Notification {
	notifs := GetNotifications(Db, User_id)
	var notifications []Notification
	for _, notif := range notifs {
		user2 := GetAccountById(Db, notif.User_id2)
		if strings.Contains(strings.ToLower(user2.Username), strings.ToLower(Search)) || strings.Contains(strings.ToLower(notif.Reason), strings.ToLower(Search)) {
			if GetAccountById(Db, notif.User_id2).PfpChanged {
				notif.PfpChanged = true
			}
			notifications = append(notifications, notif)
		}
	}
	return notifications
}

func GetNumberAllNotifsUncheck(Db *sql.DB, user_id string) int {
	row := Db.QueryRow("SELECT COUNT(*) FROM notifications WHERE user_id = ? AND checked = FALSE", user_id)
	var count int
	err := row.Scan(&count)
	if err != nil {
		fmt.Println("Error getting count notifications:", err)
	}
	return count
}
