package main

import (
	"github.com/ActiveState/tail"
	"log"
	"os/exec"
	"regexp"
)

func main() {
	var usb_disconnect = regexp.MustCompile("usb 1-2: reset full-speed USB device number 15")
	var usb_reconnect = "attach-device peaksimple /home/bohms/usb.xml"

	t, err := tail.TailFile("/var/log/syslog", tail.Config{
		Follow: true,
		ReOpen: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	for line := range t.Lines {
		log.Println(line.Text)
		if usb_disconnect.Match([]byte(line.Text)) {
			log.Println(line.Text)
			_, err := exec.Command("virsh", usb_reconnect).Output()
			if err != nil {
				log.Println(err)
			}
		}
	}
}
