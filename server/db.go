package gitweb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const Read = "Read"
const Write = "Write"

func open() *sql.DB {
	db, err := sql.Open("mysql", "gitweb:gitweb@unix(/var/run/mysql/mysql.sock)/gitweb")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Access(user, repo, access string) bool {
	db := open()
	rows, err := db.Query("SELECT a.mode, a.repo FROM gitweb.access a JOIN gitweb.user u ON a.user_id = u.id"+
		" WHERE u.name = ?"+
		" AND (a.repo = ? OR a.repo = '*')", user, repo)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var rowMode, rowRepo string
		err = rows.Scan(&rowMode, &rowRepo)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(rowMode, rowRepo, access)
		if rowMode == access || rowMode == Write {
			return true
		}
	}
	return false
}

func All(commandPath string) []string {
	db := open()
	rows, err := db.Query("SELECT u.name, u.key FROM gitweb.user u")
	if err != nil {
		log.Fatal(err)
	}

	var user, key string
	var commands []string
	for rows.Next() {
		err = rows.Scan(&user, &key)
		if err != nil {
			log.Fatal(err)
		}
		command := "command=\"" + commandPath + " " + user + "\",no-port-forwarding,no-agent-forwarding,no-X11-forwarding,no-pty " + key
		commands = append(commands, command)
	}
	return commands
}
