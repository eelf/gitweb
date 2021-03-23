package main

import (
	"github.com/eelf/gitweb"
	"github.com/keegancsmith/shell"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)


var availCmds = map[string]string{
	"git-upload-pack": gitweb.Read,
	"git-upload-archive": gitweb.Read,
	"git-receive-pack": gitweb.Write,
}

func main() {
	log.SetFlags(log.Llongfile|log.Lmicroseconds)

	if len(os.Args) != 2 {
		log.Fatal("usage: me user")
	}
	user := os.Args[1]

	origCmdStr := os.Getenv("SSH_ORIGINAL_COMMAND")

	//['/usr/local/bin/git-shell', '/usr/bin/git-shell']
	shellPath, err := exec.LookPath("git-shell")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(user, origCmdStr)

	origCmd := strings.Split(origCmdStr, " ")
	if len(origCmd) != 2 {
		log.Fatal("bad orig command:" + origCmdStr)
	}

	access, ok := availCmds[origCmd[0]]
	if !ok {
		log.Fatal("bad command")
	}
	repo := strings.Trim(origCmd[1], "'")
	repoParts := strings.Split(repo, "/")

	mapFun := func (s []string, fn func(string) string) []string {
		for i := 0; i < len(s); i++ {
			s[i] = fn(s[i])
			if len(s[i]) == 0 {
				copy(s[i:], s[i+1:])
				s = s[0:len(s) - 1]
				i--
			}
		}
		return s
	}
	repoParts = mapFun(repoParts, func(s string) string {
		return strings.Trim(s, ".")
	})
	repo = strings.Join(repoParts, "/")

	projectRoot := "/web/repos"
	fullPath := projectRoot + "/" + repo

	if fi, err := os.Stat(fullPath); err != nil || !fi.IsDir() {
		if strings.HasSuffix(fullPath, ".git") {
			log.Fatal("no such repo:" + fullPath)
		}
		fullPath += ".git"
		if fi, err = os.Stat(fullPath); err != nil || !fi.IsDir() {
			log.Fatal("no such repo2:" + fullPath)
		}
	}

	if !gitweb.Access(user, repo, access) {
		log.Fatal("no access")
	}

	log.Println("access ok")
	err = syscall.Exec(shellPath, []string{shellPath, "-c", origCmd[0] + " " + shell.EscapeArg(fullPath)}, os.Environ())
	log.Fatal(err)
}
