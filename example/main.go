package main

import "github.com/mattn/go-uv"

func main() {
	tcp, _ := uv.TcpInit()
	err := tcp.Bind("0.0.0.0", 8888)
	if err != nil {
		panic(err)
	}
	tcp.Listen(10, func(status int) {
		client, _ := tcp.Accept()
		println("server: accept")
		client.ReadStart(func(data []byte) {
			println("client: read", string(data))
			client.Write(data, func(status int) {
				println("client: written")
			})
			client.Close(func() {
				println("client: closed")
			})
			tcp.Close(func() {
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
	}()
	*/

	uv.Run()
}
