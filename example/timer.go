package main

import "github.com/mattn/go-uv"

func main() {
	timer, _ := uv.TimerInit()

	timer.Start(1000, 1000, func(status int) {
		println("timer")
		timer.Again()
	})

	uv.DefaultLoop().Run()
}
