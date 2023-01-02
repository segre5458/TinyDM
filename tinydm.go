package main

import (
	"log"
	"os"

)

func Main() {
	fTTY, err := os.OpenFile("/dev/tty7", os.O_RDWR, 0700)
	if err != nil{
		log.Fatal(err)
	}
	os.Stdout = fTTY
	os.Stdin = fTTY
	os.Stderr = fTTY

}