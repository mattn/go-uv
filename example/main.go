package main

import "github.com/mattn/go-uv"

func main() {
	tcp, _ := uv.TcpInit()
	tcp.Bind("0.0.0.0", 8888)
	tcp.Listen(10, func(status int) {
		client, _ := tcp.Accept()
		client.ReadStart(func(data []byte) {
			println("foo")
		})
	})
}
