package auth

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"

	"github.com/msteinert/pam"
	"golang.org/x/term"
)

func Login() {
	user, trans := AuthUser()
	err := trans.PutEnv("XDG_SESSION_TYPE=Xorg")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = trans.OpenSession(pam.Silent)
	if err != nil {
		log.Fatal(err.Error())
	}
	os.Setenv("HOME", user.homedir)
	os.Setenv("PWD", user.homedir)
	os.Setenv("USER", user.username)
	os.Setenv("LOGNAME", user.username)
	os.Setenv("UID", strconv.Itoa(user.uid))
	os.Setenv("XDG_CONFIG_HOME", user.homedir+"/.config")
	os.Setenv("XDG_RUNTIME_DIR", "/run/user/"+strconv.Itoa(user.uid))
	os.Setenv("XDG_SEAT", "seat0")
	os.Setenv("XDG_SESSION_CLASS", "user")
	os.Setenv("SHELL", user.loginshell)
	os.Setenv("LAGN", "en_US.utf8")
	os.Setenv("PATH", os.Getenv("PATH"))

	os.Chdir(os.Getenv("PWD"))

	
}

func AuthUser() (*sysuser, *pam.Transaction) {
	t, err := pam.StartFunc("tinydm", "", func(s pam.Style, msg string) (string, error) {
		switch s {
		case pam.PromptEchoOff:
			fmt.Print(msg)
			pw, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return "", err
			}
			fmt.Println()
			return string(pw), nil
		case pam.PromptEchoOn:
			fmt.Print(msg)
			s := bufio.NewScanner(os.Stdin)
			s.Scan()
			return s.Text(), nil
		case pam.ErrorMsg:
			fmt.Fprintf(os.Stderr, "%s\n", msg)
			return "", nil
		case pam.TextInfo:
			fmt.Println(msg)
			return "", nil
		default:
			return "", errors.New("unrecognized message style")
		}
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "start: %s\n", err.Error())
		os.Exit(1)
	}
	err = t.Authenticate(pam.Silent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "authenticate: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("authentication succeeded!")

	pamUsr, err := t.GetItem(pam.User)
	if err != nil {
		log.Fatal(err.Error())
	}
	usr, _ := user.Lookup(pamUsr)
	return getsysUser(usr), t
}
