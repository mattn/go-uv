package main

import (
	"github.com/mattn/go-uv"
	"log"
)

func main() {
	timer, _ := uv.TimerInit(nil)
	timer.Start(func(h *uv.Handle, status int) {
		log.Println("tick")
	}, 1000, 1000)

	uv.DefaultLoop().Run()
}
