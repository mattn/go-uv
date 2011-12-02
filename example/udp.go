package main

import "github.com/mattn/go-uv"

func main() {
	udp, _ := uv.UdpInit()
	err := udp.Bind(uv.Ip4Addr("0.0.0.0", 8888), 0)
	if err != nil {
		panic(err)
	}
	udp.Listen(10, func(status int) {
		client, _ := udp.Accept()
		println("server: accept")
		client.RecvStart(func(data []byte, sa uv.SockaddrIn, flags uint) {
			println("client: read", string(data))
			client.Send(data, sa, func(status int) {
				println("client: written")
			})
			client.Close(func() {
				println("client: closed")
			})
			udp.Close(func() {
				println("server: closed")
			})
		})
	})

	/*
		go func() {
			time.Sleep(1e9)
			tcp, _ := uv.TcpInit()
			tcp.Connect("0.0.0.0", 8888, func(status int) {
				println(status)
				tcp.Write([]byte("Hello World"), func(status int) {
					println("sender: sent!", status)
				})
			})
			uv.DefaultLoop().Run()
		}()
	*/

	uv.DefaultLoop().Run()
}
