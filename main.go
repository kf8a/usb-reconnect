package main

import (
	"fmt"
	"github.com/ActiveState/tail"
	"regexp"
)

func main() {
	var usb_disconnect = regexp.MustCompile("usb 1-2: reset full-speed USB device number 15")

	t, _ := tail.TailFile("/var/log/syslog", tail.Config{Follow: true})
	for line := range t.Lines {
		if usb_disconnect.Match([]byte(line.Text)) {
			fmt.Println(line.Text)
		}
	}
}
