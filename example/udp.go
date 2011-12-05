package main

import "github.com/mattn/go-uv"

func main() {
	udp, _ := uv.UdpInit(nil)
	err := udp.Bind(uv.Ip4Addr("0.0.0.0", 8888), 0)
	if err != nil {
		panic(err)
	}
	println("udp: start")
	udp.RecvStart(func(h *uv.Handle, data []byte, sa uv.SockaddrIn, flags uint) {
		println("udp: read", string(data))
		udp.Send(data, sa, func(r *uv.Request, status int) {
			println("udp: written")
		})
		udp.Close(func(h *uv.Handle) {
			println("udp: closed")
		})
	})

	uv.DefaultLoop().Run()
}
