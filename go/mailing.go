package forum

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"html/template"
	"log"
	"net/smtp"
)

func SendPasswordResetEmail(email, token string) error {
	route := "http://localhost:8080/reset-password?token=" + token + "/" // A CHANGER

	var body bytes.Buffer
	t, _ := template.ParseFiles("./go/template_resetpassword.html")
	t.Execute(&body, struct{ Route string }{Route: route})

	auth := smtp.PlainAuth("", "forumprojetynov@gmail.com", "ljhu jgfl atnq lkbh", "smtp.gmail.com")

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := "Subject: Réinisialisation de mot de passe\n" + headers + "\n\n" + body.String()

	err := smtp.SendMail("smtp.gmail.com:587", auth, "forumprojetynov@no-reply.com", []string{email}, []byte(msg))

	if err != nil {
		return err
	}
	return nil
}

func SendCreatedAccountEmail(email, username string) error {
	var body bytes.Buffer
	t, _ := template.ParseFiles("./go/template_accountcreated.html")
	t.Execute(&body, struct{ Username string }{Username: username})

	auth := smtp.PlainAuth("", "forumprojetynov@gmail.com", "ljhu jgfl atnq lkbh", "smtp.gmail.com")

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := "Subject: Bienvenue sur ParlonsSanté !\n" + headers + "\n\n" + body.String()

	err := smtp.SendMail("smtp.gmail.com:587", auth, "forumprojetynov@no-reply.com", []string{email}, []byte(msg))

	if err != nil {
		return err
	}
	return nil
}

func EncodeToken(email string) string {
	token := base64.URLEncoding.EncodeToString([]byte(email))
	return token
}

func DecodeToken(token string) (string, error) {
	emailBytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}
	email := string(emailBytes)
	return email, nil
}

func InvalidAllMail(db *sql.DB) {
	mails, err := GetAllMail(db)
	if err != nil {
		log.Println(err)
	}
	for _, mail := range mails {
		ResetPasswordMap[EncodeToken(mail)+"/"] = "invalid"
	}
}
