package main

import "github.com/mattn/go-uv"

func main() {
	timer, _ := uv.TimerInit(nil)
	timer.Start(func(h *uv.Handle, status int) {
		println("timer")
	}, 1000, 1000)

	uv.DefaultLoop().Run()
}
