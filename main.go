package main

import (
	"github.com/ActiveState/tail"
	"log"
	"os/exec"
	"regexp"
)

func main() {
	var usb_found = regexp.MustCompile("usb 1-2: New USB device found, idVendor=0724, idProduct=0004")
	var usb_disconnect = regexp.MustCompile("usb 1-2: reset full-speed USB device")
	var usb_message = regexp.MustCompile("usb")
	var usb_reconnect = "attach-device peaksimple /home/bohms/usb.xml"

	t, err := tail.TailFile("/var/log/syslog", tail.Config{
		Follow: true,
		ReOpen: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	for line := range t.Lines {
		if usb_message.Match([]byte(line.Text)) {
			log.Println(line.Text)
		}
		if usb_disconnect.Match([]byte(line.Text)) || usb_found.Match([]byte(line.Text)) {
			log.Println("usb disconnect detected ... reconnecting")
			log.Println(line.Text)
			output, err := exec.Command("virsh", usb_reconnect).Output()
			log.Println(string(output))
			if err != nil {
				log.Println(err)
			}
		}
	}
}
