package auth

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/msteinert/pam"
	"golang.org/x/term"
)

func Login() {
	user, trans := AuthUser()

	var desktop []fs.FileInfo
	files, err := ioutil.ReadDir("/usr/share/xsessions/")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".desktop" {
			desktop = append(desktop, file)
		}
	}
	if len(desktop) == 0 {
		log.Fatal("No Desktop Entry")
	}
	fmt.Printf("\n")
	for i, w := range desktop {
		fmt.Printf("[%d] %s", i, w.Name())
		if i != len(desktop)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Printf("\nSelect [0]: ")
	selection := bufio.NewScanner(os.Stdin)
	selection.Scan()
	id, err := strconv.Atoi(selection.Text())
	if err != nil {
		log.Fatal(err.Error())
	}
	if id >= len(desktop) || id <= 0 {
		log.Fatal("Not Found Such Desktop Entry")
	}
	wm := desktop[id]

	var execCmd *exec.Cmd
	fp, err := os.Open("/usr/share/xsessions/"+wm.Name())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		slice := strings.Split(scanner.Text(), "=")
		if slice[0] == "Exec" {
			execCmd = exec.Command(slice[1])
		}
	}
	if execCmd == nil {
		log.Fatal("Not Found Exec Comman in Desktop Entry")
	}

	err = execCmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = execCmd.Wait()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = trans.PutEnv("XDG_SESSION_TYPE=Xorg")
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

	d := &dbus{}
	d.launch(user)
	d.intrrupt()
	if trans != nil {
		if err := trans.SetCred(pam.DeleteCred); err != nil {
			log.Fatal(err.Error())
		}
		if err := trans.CloseSession(pam.Silent); err != nil {
			log.Fatal(err.Error())
		}
		trans = nil
	}
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
