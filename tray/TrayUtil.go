package tray

import (
	"github.com/getlantern/systray"
)

var data []byte

func InitTray(onExit func(), icoData []byte) {
	data = icoData
	systray.Run(onReady, onExit)
}

func onReady() {
	//data, err := ioutil.ReadFile("favicon.ico")
	//if err != nil {
	//	fmt.Println("Error reading ICO file:", err)
	//}
	systray.SetIcon(data)
	systray.SetTooltip("MiHome")

	mQuit := systray.AddMenuItem("退出", "")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}
