package auth

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/msteinert/pam"
)

func Login() {
	_ = AuthUser()
}

func AuthUser() *sysuser {
	trans, _ := pam.StartFunc("tidydm", "segre", func(s pam.Style, msg string) (string, error) {
		hostname, _ := os.Hostname()
		fmt.Printf("%s login: \n", hostname)
		fmt.Print("Password: ")

		// fd := os.Stdout.Fd()
		// c := make(chan os.Signal, 1)
		// signal.Notify(c)

		// go handleInt
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return "", err
		}
		fmt.Println()
		return input[:len(input)-1], nil
	})
	_ = trans.Authenticate(pam.Silent)
	log.Print("Authenticate OK")
	return nil
}
