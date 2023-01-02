package auth

import (
	"log"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

type sysuser struct {
	username   string
	uid        int
	gid        int
	homedir    string
	loginshell string
}

func getsysUser(usr *user.User) *sysuser {
	var res sysuser
	res.username = usr.Username
	res.uid, _ = strconv.Atoi(usr.Uid)
	res.gid, _ = strconv.Atoi(usr.Gid)
	res.homedir = usr.HomeDir
	res.loginshell = getShell(res.uid)

	return &res
}

func getShell(uid int) string {
	out, err := exec.Command("/usr/bin/getent", "passwd", strconv.Itoa(uid)).Output()
	if err != nil{
		log.Fatal(err.Error())
	}

	ent := strings.Split(strings.TrimSuffix(string(out), "\n"), ":")
	return ent[6]
}
