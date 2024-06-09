package main

import (
	"bemfa-demo/mqtt"
	"bemfa-demo/tray"
	"embed"
	"io/fs"
)

//go:embed favicon.ico
var iconFS embed.FS

func main() {
	go mqtt.Launch()
	file, _ := fs.ReadFile(iconFS, "favicon.ico")
	ints := make(chan int, 1)
	tray.InitTray(func() {
		ints <- 0
	}, file)
	<-ints
}
