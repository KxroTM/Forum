package forum

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/smtp"
)

func SendEmail(email, token string) {

	route := "http://localhost:8070/reset-password?token=" + token + "/" // A CHANGER

	var body bytes.Buffer
	t, _ := template.ParseFiles("./app/template.html")
	t.Execute(&body, struct{ Route string }{Route: route})

	auth := smtp.PlainAuth("", "forumprojetynov@gmail.com", "ljhu jgfl atnq lkbh", "smtp.gmail.com")

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := "Subject: Réinisialisation de mot de passe\n" + headers + "\n\n" + body.String()

	err := smtp.SendMail("smtp.gmail.com:587", auth, "forumprojetynov@no-reply.com", []string{email}, []byte(msg))

	if err != nil {
		fmt.Println(err)
	}
}

func EncodeToken(email string) string {
	token := base64.URLEncoding.EncodeToString([]byte(email))
	return token
}

func DecodeToken(token string) string {
	emailBytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return ""
	}
	email := string(emailBytes)
	return email
}

func InvalidAllMail() {
	mails := GetAllMail()
	for _, mail := range mails {
		ResetPasswordMap[EncodeToken(mail)+"/"] = "invalid"
	}
}
