package main

import (
	"github.com/mattn/go-uv"
	"log"
)

func main() {
	udp, _ := uv.UdpInit(nil)
	addr, err := uv.Ip4Addr("0.0.0.0", 8888)
	if err != nil {
		log.Fatal(err)
	}
	err = udp.Bind(addr, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("udp: start")
	udp.RecvStart(func(h *uv.Handle, data []byte, sa uv.SockaddrIn, flags uint) {
		log.Println("udp: read", string(data))
		udp.Send(data, sa, func(r *uv.Request, status int) {
			log.Println("udp: written")
		})
		udp.Close(func(h *uv.Handle) {
			log.Println("udp: closed")
		})
	})

	uv.DefaultLoop().Run()
}
