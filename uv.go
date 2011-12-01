package uv

/*
#include <stdlib.h>
#include <uv/uv.h>

extern void __uv_connect_cb(void* p, int status);
static void _uv_connect_cb(uv_stream_t* server, int status) {
	__uv_connect_cb(server, status);
}

static void _uv_read_cb(uv_stream_t* stream, ssize_t nread, uv_buf_t buf) {
}

static uv_buf_t _uv_alloc_cb(uv_handle_t* handle, size_t suggested_size) {
    char* buf;
    buf = (char*)malloc(suggested_size);
    return uv_buf_init(buf, suggested_size);
}

static int _uv_listen(uv_stream_t* stream, int backlog) {
	return uv_listen(stream, backlog, _uv_connect_cb);
}

static int _uv_read_start(uv_stream_t* stream) {
	return uv_read_start(stream, _uv_alloc_cb, _uv_read_cb);
}

static uv_stream_t* _uv_tcp_to_string(uv_tcp_t* tcp) {
	return (uv_stream_t*) tcp;
}

#cgo CFLAGS: -static
#cgo LDFLAGS: -static -luv -lws2_32
*/
import "C"
import "errors"
import "unsafe"

type Tcp struct {
	t C.uv_tcp_t
}

type Stream struct {
	s C.uv_stream_t
}

type Client struct {
	c C.uv_stream_t
}

type callback_info struct {
	connect_cb func(int)
	read_cb func([]byte)
}

func TcpInit() (tcp *Tcp, err error) {
	var t C.uv_tcp_t

	r := C.uv_tcp_init(C.uv_default_loop(), &t)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return nil, errors.New(C.GoString(C.uv_strerror(e)))
	}
	t.data = unsafe.Pointer(&callback_info{})
	return &Tcp {t}, nil
}

func (tcp *Tcp) Bind(host string, port uint16) (err error) {
	phost := C.CString(host)
	defer C.free(unsafe.Pointer(phost))
	r := C.uv_tcp_bind(&tcp.t, C.uv_ip4_addr(phost, C.int(port)));
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (tcp *Tcp) Listen(backlog int, cb func(int)) (err error) {
	cbi := (*callback_info)(tcp.t.data)
	cbi.connect_cb = cb
	r := C._uv_listen(C._uv_tcp_to_string(&tcp.t), C.int(backlog))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (tcp *Tcp) Accept() (client *Client, err error) {
	var c C.uv_stream_t
	r := C.uv_accept(C._uv_tcp_to_string(&tcp.t), &c)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return nil, errors.New(C.GoString(C.uv_strerror(e)))
	}
	return &Client{c}, nil
}

func (client *Client) ReadStart(cb func([]byte)) (err error) {
	cbi := (*callback_info)(client.c.data)
	cbi.read_cb = cb
	r := C._uv_read_start(&client.c)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

//export __uv_connect_cb
func __uv_connect_cb(p unsafe.Pointer, status int) {
	cbi := (*callback_info)(p)
	cbi.connect_cb(status)
}
