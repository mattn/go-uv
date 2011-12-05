package main

import "github.com/mattn/go-uv"

func main() {
	tcp, _ := uv.TcpInit(nil)
	err := tcp.Bind(uv.Ip4Addr("0.0.0.0", 8888))
	if err != nil {
		panic(err)
	}
	tcp.Listen(10, func(h *uv.Handle, status int) {
		client, _ := tcp.Accept()
		println("server: accept")
		client.ReadStart(func(h *uv.Handle, data []byte) {
			println("client: read", string(data))
			client.Write(data, func(r *uv.Request, status int) {
				println("client: written")
			})
			client.Close(func(h *uv.Handle) {
				println("client: closed")
			})
			tcp.Close(func(h *uv.Handle) {
				println("server: closed")
			})
		})
	})

	uv.DefaultLoop().Run()
}
