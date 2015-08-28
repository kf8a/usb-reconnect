package main

import (
	"github.com/ActiveState/tail"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"time"
)

var (
	usbFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "usb_errors_total",
		Help: "Number of usb errors.",
	})
	usbNewEvents = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "usb_new_total",
		Help: "Number of new usb devices found.",
	})
)

func init() {
	prometheus.MustRegister(usbFailures)
	prometheus.MustRegister(usbNewEvents)
}

func findUsbErrors() {
	var rundelayTime = time.Now().Add(1 * time.Minute)
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
		if time.Now().After(rundelayTime) {
			if usb_message.Match([]byte(line.Text)) {
				log.Println(line.Text)
			}
			if usb_disconnect.Match([]byte(line.Text)) {
				usbFailures.Inc()
				log.Println("usb disconnected ...")
			}
			if usb_found.Match([]byte(line.Text)) {
				usbNewEvents.Inc()
				log.Println("usb disconnect detected ... reconnecting")
				time.Sleep(600)
				output, err := exec.Command("virsh", usb_reconnect).Output()
				log.Println(string(output))
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func main() {

	go findUsbErrors()

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(":9101", nil)

}
