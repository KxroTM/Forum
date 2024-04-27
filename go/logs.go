package forum

import (
	"os"
	"time"
)

func AccountLog(connection string) {
	filePath := "./logs/account_connection.log"
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}

	currentDate := time.Now().In(time.FixedZone("CET", 1*60*60)).Format("02/01/2006")
	currentHour := time.Now().In(time.FixedZone("CET", 1*60*60)).Format("15:04:05")
	_, err = file.WriteString("[" + currentDate + "] " + "[" + currentHour + "]" + " : " + connection + "\n")
	if err != nil {
		return
	}
	defer file.Close()
}

func IPsLog(connection string) {
	filePath := "./logs/ip_connection.log"
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}

	currentDate := time.Now().In(time.FixedZone("CET", 1*60*60)).Format("02/01/2006")
	currentHour := time.Now().In(time.FixedZone("CET", 1*60*60)).Format("15:04:05")
	_, err = file.WriteString("[" + currentDate + "] " + "[" + currentHour + "]" + " : " + connection + "\n")
	if err != nil {
		return
	}
	defer file.Close()
}
