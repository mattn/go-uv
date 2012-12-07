package main

import (
	"fmt"
	"github.com/mattn/go-uv"
	"github.com/mattn/go-gtk/gtk"
	"time"
)

func main() {
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Clock")
	vbox := gtk.NewVBox(false, 1)
	label := gtk.NewLabel("")
	vbox.Add(label)
	window.Add(vbox)
	window.SetDefaultSize(300, 20)
	window.ShowAll()

	timer, _ := uv.TimerInit(nil)
	timer.Start(func(h *uv.Handle, status int) {
		label.SetLabel(fmt.Sprintf("%v", time.Now()))
	}, 1000, 1000)

	idle, _ := uv.IdleInit(nil)
	idle.Start(func(h *uv.Handle, status int) {
		gtk.MainIterationDo(false)
	})

	window.Connect("destroy", func() {
		timer.Close(nil)
		idle.Close(nil)
	})

	uv.DefaultLoop().Run()
}
