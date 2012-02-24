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
		line := ""
		client.ReadStart(func(h *uv.Handle, data []byte) {
			if data == nil {
				println("client: closed")
				client.Close(nil)
				return
			}
			s := string(data)
			print(s)
			line += s
			if s[len(s)-1] == '\n' {
				client.Write([]byte(line), func(r *uv.Request, status int) {
					println("client: written")
				})
			}
		})
	})

	uv.DefaultLoop().Run()
}
