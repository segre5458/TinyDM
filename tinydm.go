package main

import (
	// "log"
	// "os"

	auth "github.com/segre5458/tinyDM/auth"
)

func Main() {
	// fTTY, err := os.OpenFile("/dev/tty7", os.O_RDWR, 0700)
	// if err != nil{
	// 	log.Fatal(err)
	// }
	// os.Stdout = fTTY
	// os.Stdin = fTTY
	// os.Stderr = fTTY
	auth.Login()
	// fTTY.Close()
}