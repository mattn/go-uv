package uv

/*
#include <stdlib.h>
#include <uv/uv.h>

extern void __uv_connection_cb(void* p, int status);
static void _uv_connection_cb(uv_stream_t* server, int status) {
	__uv_connection_cb(server->data, status);
}

extern void __uv_udp_recv_cb(void* p, int nread, void* buf, void* sa, unsigned flags);
static void _uv_udp_recv_cb(uv_udp_t* udp, ssize_t nread, uv_buf_t buf, struct sockaddr* addr, unsigned flags) {
	__uv_udp_recv_cb(udp->data, nread, buf.base, addr, flags);
}

extern void __uv_udp_send_cb(void* p, int status);
static void _uv_udp_send_cb(uv_udp_send_t* req, int status) {
	__uv_udp_send_cb(req->handle->data, status);
}

static int _uv_udp_send(uv_udp_send_t* req, uv_udp_t* handle, uv_buf_t bufs[], int bufcnt, struct sockaddr_in addr) {
	return uv_udp_send(req, handle, bufs, bufcnt, addr, _uv_udp_send_cb);
}

static int _uv_udp_send6(uv_udp_send_t* req, uv_udp_t* handle, uv_buf_t bufs[], int bufcnt, struct sockaddr_in6 addr) {
	return uv_udp_send6(req, handle, bufs, bufcnt, addr, _uv_udp_send_cb);
}

static uv_buf_t _uv_alloc_cb(uv_handle_t* handle, size_t suggested_size) {
    char* buf;
    buf = (char*)malloc(suggested_size);
    return uv_buf_init(buf, suggested_size);
}

extern void __uv_shutdown_cb(void* p, int status);
static void _uv_shutdown_cb(uv_shutdown_t* req, int status) {
	__uv_shutdown_cb(req->handle->data, status);
}

static int _uv_shutdown(uv_shutdown_t* req, uv_stream_t* handle) {
	return uv_shutdown(req, handle, _uv_shutdown_cb);
}


extern void __uv_close_cb(void* p);
static void _uv_close_cb(uv_handle_t* handle) {
	__uv_close_cb(handle->data);
}

static int _uv_udp_recv_start(uv_udp_t* udp) {
	return uv_udp_recv_start(udp, _uv_alloc_cb, _uv_udp_recv_cb);
}

static int _uv_listen(uv_stream_t* stream, int backlog) {
	return uv_listen(stream, backlog, _uv_connection_cb);
}

static void _uv_close(uv_handle_t* handle) {
	uv_close(handle, _uv_close_cb);
}

static uv_stream_t* _uv_udp_to_stream(uv_udp_t* udp) {
	return (uv_stream_t*) udp;
}

static uv_handle_t* _uv_udp_to_handle(uv_udp_t* udp) {
	return (uv_handle_t*) udp;
}
*/
import "C"
import "errors"
import "unsafe"

type Udp struct {
	t *C.uv_udp_t
}

type udp_callback_info struct {
	connection_cb func(int)
	connect_cb    func(int)
	udp_recv_cb       func([]byte, SockaddrIn, uint)
	send_cb      func(int)
	close_cb      func()
	shutdown_cb    func(int)
}

func UdpInit() (udp *Udp, err error) {
	var t C.uv_udp_t

	r := C.uv_udp_init(C.uv_default_loop(), &t)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return nil, errors.New(C.GoString(C.uv_strerror(e)))
	}
	t.data = unsafe.Pointer(&udp_callback_info{})
	return &Udp{&t}, nil
}

func (udp *Udp) Bind(sa SockaddrIn, flags uint) (err error) {
	var r C.int
	sa4, is_v4 := sa.(*SockaddrIn4)
	if is_v4 {
		r = C.uv_udp_bind(udp.t, sa4.sa, C.uint(flags))
	} else {
		sa6, is_v6 := sa.(*SockaddrIn6)
		if is_v6 {
			r = C.uv_udp_bind6(udp.t, sa6.sa, C.uint(flags))
		}
	}
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (udp *Udp) Listen(backlog int, cb func(int)) (err error) {
	cbi := (*udp_callback_info)(udp.t.data)
	cbi.connection_cb = cb
	r := C._uv_listen(C._uv_udp_to_stream(udp.t), C.int(backlog))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (udp *Udp) Accept() (client *Udp, err error) {
	c, err := UdpInit()
	if err != nil {
		return nil, err
	}
	r := C.uv_accept(C._uv_udp_to_stream(udp.t), C._uv_udp_to_stream(c.t))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return nil, errors.New(C.GoString(C.uv_strerror(e)))
	}
	return &Udp{c.t}, nil
}

func (udp *Udp) RecvStart(cb func([]byte, SockaddrIn, uint)) (err error) {
	cbi := (*udp_callback_info)(udp.t.data)
	cbi.udp_recv_cb = cb
	r := C._uv_udp_recv_start(udp.t)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (udp *Udp) RecvStop() (err error) {
	r := C.uv_udp_recv_stop(udp.t)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (udp *Udp) Send(b []byte, sa SockaddrIn, cb func(int)) (err error) {
	cbi := (*udp_callback_info)(udp.t.data)
	cbi.send_cb = cb
	var req C.uv_udp_send_t
	buf := C.uv_buf_init((*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)))
	var r C.int
	sa4, is_v4 := sa.(*SockaddrIn4)
	if is_v4 {
		r = C._uv_udp_send(&req, udp.t, &buf, 1, sa4.sa)
	} else {
		sa6, is_v6 := sa.(*SockaddrIn6)
		if is_v6 {
			r = C._uv_udp_send6(&req, udp.t, &buf, 1, sa6.sa)
		}
	}
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (udp *Udp) Shutdown(cb func(int)) {
	cbi := (*udp_callback_info)(udp.t.data)
	cbi.shutdown_cb = cb
	var req C.uv_shutdown_t
	C._uv_shutdown(&req, C._uv_udp_to_stream(udp.t))
}

func (udp *Udp) Close(cb func()) {
	cbi := (*udp_callback_info)(udp.t.data)
	cbi.close_cb = cb
	C._uv_close(C._uv_udp_to_handle(udp.t))
}

func (udp *Udp) IsActive() bool {
	if C.uv_is_active(C._uv_udp_to_handle(udp.t)) != 0 {
		return true
	}
	return false
}

//export __uv_udp_recv_cb
func __uv_udp_recv_cb(p unsafe.Pointer, nread int, buf unsafe.Pointer, sa unsafe.Pointer, flags uint) {
	cbi := (*udp_callback_info)(p)
	if cbi.udp_recv_cb != nil {
		psa := &SockaddrIn4{*(*C.struct_sockaddr_in)(sa)}
		cbi.udp_recv_cb((*[1 << 30]byte)(unsafe.Pointer(buf))[0:nread], psa, flags)
	}
}

//export __uv_udp_send_cb
func __uv_udp_send_cb(p unsafe.Pointer, status int) {
	cbi := (*udp_callback_info)(p)
	if cbi.send_cb != nil {
		cbi.send_cb(status)
	}
}
