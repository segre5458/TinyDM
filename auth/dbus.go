package auth

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type dbus struct {
	pid     int
	address string
}

func (d *dbus) launch(usr *sysuser) {
	outTxt, err := exec.Command("dbus-launch").Output()
	if err != nil {
		log.Fatal(err.Error())
	}
	if outTxt == nil {
		log.Fatal("D-Bus Not Respond")
	}
	scanner := bufio.NewScanner(strings.NewReader(string(outTxt)))
	for scanner.Scan() {
		slice := strings.Split(scanner.Text(), "=")
		key := slice[0]
		value := slice[1]
		switch key {
		case "DBUS_SESSION_BUS_ADDRESS":
			d.address = value
			os.Setenv("DBUS_SESSION_BUS_ADDRESS", value)
		case "DBUS_SESSION_BUS_PID":
			d.pid, err = strconv.Atoi(value)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}
}

func (d *dbus) intrrupt(){
	if d.pid <= 0{
		return
	}
	proc,err := os.FindProcess(d.pid)
	if err != nil{
		log.Fatal(err.Error())
	}
	if proc != nil{
		err = proc.Signal(os.Interrupt)
		if err != nil{
			log.Fatal(err.Error())
		}
	}
}
