package uv

/*
#include <uv/uv.h>
*/
import "C"
import "unsafe"

type Udp struct {
	u *C.uv_udp_t
	l *C.uv_loop_t
	Handle
}

func UdpInit(loop *Loop) (udp *Udp, err error) {
	var u C.uv_udp_t

	if loop == nil {
		loop = DefaultLoop()
	}
	r := C.uv_udp_init(loop.l, &u)
	if r != 0 {
		return nil, udp.GetLoop().LastError().Error()
	}
	u.data = unsafe.Pointer(&callback_info{})
	return &Udp{&u, loop.l, Handle{(*C.uv_handle_t)(unsafe.Pointer(&u)), u.data}}, nil
}

func (udp *Udp) GetLoop() *Loop {
	return &Loop{udp.l}
}

func (udp *Udp) Bind(sa SockaddrIn, flags uint) (err error) {
	var r int
	sa4, is_v4 := sa.(*SockaddrIn4)
	if is_v4 {
		r = int(C.uv_udp_bind(udp.u, sa4.sa, C.uint(flags)))
	} else {
		sa6, is_v6 := sa.(*SockaddrIn6)
		if is_v6 {
			r = int(C.uv_udp_bind6(udp.u, sa6.sa, C.uint(flags)))
		}
	}
	if r != 0 {
		return udp.GetLoop().LastError().Error()
	}
	return nil
}

func (udp *Udp) RecvStart(cb func(*Handle, []byte, SockaddrIn, uint)) (err error) {
	cbi := (*callback_info)(udp.u.data)
	cbi.udp_recv_cb = cb
	r := uv_udp_recv_start(udp.u)
	if r != 0 {
		return udp.GetLoop().LastError().Error()
	}
	return nil
}

func (udp *Udp) RecvStop() (err error) {
	r := uv_udp_recv_stop(udp.u)
	if r != 0 {
		return udp.GetLoop().LastError().Error()
	}
	return nil
}

func (udp *Udp) Send(b []byte, sa SockaddrIn, cb func(*Request, int)) (err error) {
	cbi := (*callback_info)(udp.u.data)
	cbi.udp_send_cb = cb
	buf := C.uv_buf_init((*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)))
	var r int
	sa4, is_v4 := sa.(*SockaddrIn4)
	if is_v4 {
		r = uv_udp_send(udp.u, &buf, 1, sa4.sa)
	} else {
		sa6, is_v6 := sa.(*SockaddrIn6)
		if is_v6 {
			r = uv_udp_send6(udp.u, &buf, 1, sa6.sa)
		}
	}
	if r != 0 {
		return udp.GetLoop().LastError().Error()
	}
	return nil
}

func (udp *Udp) Shutdown(cb func(*Request, int)) {
	cbi := (*callback_info)(udp.u.data)
	cbi.shutdown_cb = cb
	uv_shutdown((*C.uv_stream_t)(unsafe.Pointer(udp.u)))
}

func (udp *Udp) GetSockname() (sa *Sockaddr, err error) {
	var csa C.struct_sockaddr
	r := uv_udp_getsockname(udp.u, &csa)
	if r != 0 {
		return nil, udp.GetLoop().LastError().Error()
	}
	return &Sockaddr{csa}, nil
}
