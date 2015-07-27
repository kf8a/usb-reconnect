package main

import (
	"fmt"
	"github.com/ActiveState/tail"
)

func main() {
	t, _ := tail.TailFile("/var/log/nginx.log", tail.Config{Follow: true})
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}
