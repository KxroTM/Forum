// Fichier main.go
package main

import (
	"database/sql"
	"fmt"

	forum "forum/go"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func main() {

}

func init() {
	var err error
	Db, err = sql.Open("sqlite3", "./database/db.db")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture de la base de données:", err)
		return
	}
	forum.UpdateDb(Db)
}
