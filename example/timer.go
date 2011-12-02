package main

import "github.com/mattn/go-uv"

func main() {
	timer, _ := uv.TimerInit(nil)
	timer.Start(1000, 1000, func(status int) {
		println("timer")
	})

	uv.DefaultLoop().Run()
}
