package uv

/*
#include <uv.h>
*/
import "C"
import "unsafe"

type Tcp struct {
	t *C.uv_tcp_t
	l *C.uv_loop_t
	Handle
}

func TcpInit(loop *Loop) (tcp *Tcp, err error) {
	var t C.uv_tcp_t

	if loop == nil {
		loop = DefaultLoop()
	}
	r := C.uv_tcp_init(loop.l, &t)
	if r != 0 {
		return nil, tcp.GetLoop().LastError().Error()
	}
	t.data = unsafe.Pointer(&callback_info{})
	return &Tcp{&t, loop.l, Handle{(*C.uv_handle_t)(unsafe.Pointer(&t)), t.data}}, nil
}

func (tcp *Tcp) GetLoop() *Loop {
	return &Loop{tcp.l}
}

func (tcp *Tcp) Bind(sa SockaddrIn) (err error) {
	var r int
	sa4, is_v4 := sa.(*SockaddrIn4)
	if is_v4 {
		r = uv_tcp_bind(tcp.t, sa4.sa)
	} else {
		sa6, is_v6 := sa.(*SockaddrIn6)
		if is_v6 {
			r = uv_tcp_bind6(tcp.t, sa6.sa)
		}
	}
	if r != 0 {
		return tcp.GetLoop().LastError().Error()
	}
	return nil
}

func (tcp *Tcp) Nodelay(enable bool) (err error) {
	r := uv_tcp_nodelay(tcp.t, enable)
	if r != 0 {
		return tcp.GetLoop().LastError().Error()
	}
	return nil
}

func (tcp *Tcp) Keepalive(enable bool, delay uint) (err error) {
	r := uv_tcp_keepalive(tcp.t, enable, delay)
	if r != 0 {
		return tcp.GetLoop().LastError().Error()
	}
	return nil
}

func (tcp *Tcp) SimultaneousAccepts(enable bool) (err error) {
	var v C.int
	if enable {
		v = 1
	}
	r := C.uv_tcp_simultaneous_accepts(tcp.t, v)
	if r != 0 {
		return tcp.GetLoop().LastError().Error()
	}
	return nil
}

func (tcp *Tcp) Connect(sa SockaddrIn, cb func(*Request, int)) (err error) {
	cbi := (*callback_info)(tcp.t.data)
	cbi.connect_cb = cb
	var r int
	sa4, is_v4 := sa.(*SockaddrIn4)
	if is_v4 {
		r = uv_tcp_connect(tcp.t, sa4.sa)
	} else {
		sa6, is_v6 := sa.(*SockaddrIn6)
		if is_v6 {
			r = uv_tcp_connect6(tcp.t, sa6.sa)
		}
	}
	if r != 0 {
		return tcp.GetLoop().LastError().Error()
	}
	return nil
}

func (tcp *Tcp) Listen(backlog int, cb func(*Handle, int)) (err error) {
	cbi := (*callback_info)(tcp.t.data)
	cbi.connection_cb = cb
	r := uv_listen((*C.uv_stream_t)(unsafe.Pointer(tcp.t)), backlog)
	if r != 0 {
		return tcp.GetLoop().LastError().Error()
	}
	return nil
}

func (tcp *Tcp) Accept() (client *Tcp, err error) {
	c, err := TcpInit(tcp.GetLoop())
	if err != nil {
		return nil, err
	}
	r := uv_accept((*C.uv_stream_t)(unsafe.Pointer(tcp.t)), (*C.uv_stream_t)(unsafe.Pointer(c.t)))
	if r != 0 {
		return nil, tcp.GetLoop().LastError().Error()
	}
	return &Tcp{c.t, tcp.l, Handle{(*C.uv_handle_t)(unsafe.Pointer(c.t)), c.t.data}}, nil
}

func (tcp *Tcp) ReadStart(cb func(*Handle, []byte)) (err error) {
	cbi := (*callback_info)(tcp.t.data)
	cbi.read_cb = cb
	r := uv_read_start((*C.uv_stream_t)(unsafe.Pointer(tcp.t)))
	if r != 0 {
		return tcp.GetLoop().LastError().Error()
	}
	return nil
}

func (tcp *Tcp) ReadStop() (err error) {
	r := uv_read_stop((*C.uv_stream_t)(unsafe.Pointer(tcp.t)))
	if r != 0 {
		return tcp.GetLoop().LastError().Error()
	}
	return nil
}

func (tcp *Tcp) Write(b []byte, cb func(*Request, int)) (err error) {
	cbi := (*callback_info)(tcp.t.data)
	cbi.write_cb = cb
	buf := uv_buf_init(b)
	r := uv_write((*C.uv_stream_t)(unsafe.Pointer(tcp.t)), &buf, 1)
	if r != 0 {
		return tcp.GetLoop().LastError().Error()
	}
	return nil
}

func (tcp *Tcp) Shutdown(cb func(*Request, int)) (err error) {
	cbi := (*callback_info)(tcp.t.data)
	cbi.shutdown_cb = cb
	r := uv_shutdown((*C.uv_stream_t)(unsafe.Pointer(tcp.t)))
	if r != 0 {
		return tcp.GetLoop().LastError().Error()
	}
	return nil
}

func (tcp *Tcp) GetSockname() (sa *Sockaddr, err error) {
	var csa C.struct_sockaddr
	r := uv_tcp_getsockname(tcp.t, &csa)
	if r != 0 {
		return nil, tcp.GetLoop().LastError().Error()
	}
	return &Sockaddr{csa}, nil
}

func (tcp *Tcp) GetPeername() (sa *Sockaddr, err error) {
	var csa C.struct_sockaddr
	r := uv_tcp_getpeername(tcp.t, &csa)
	if r != 0 {
		return nil, tcp.GetLoop().LastError().Error()
	}
	return &Sockaddr{csa}, nil
}
