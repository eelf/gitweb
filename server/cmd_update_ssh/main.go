package main

import (
	"github.com/eelf/gitweb"
	"log"
	"os"
	"strings"
)

func main() {
	commandPath := "/root/go/bin/cmd_gitshell"

	commands := gitweb.All(commandPath)
	authKeys := strings.Join(commands, "\n")

	f, err := os.Create("/web/repos/.ssh/authorized_keys")
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.WriteString(authKeys)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ok")
}
