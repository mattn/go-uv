package uv

/*
#include <stdlib.h>
#include <uv/uv.h>

extern void __uv_connect_cb(uv_connect_t* req, int status);
extern void __uv_connection_cb(uv_stream_t* stream, int status);
extern void __uv_write_cb(uv_write_t* req, int status);
extern void __uv_read_cb(uv_stream_t* stream, ssize_t nread, uv_buf_t buf);
extern void __uv_udp_recv_cb(uv_udp_t* handle, ssize_t nread, uv_buf_t buf, struct sockaddr* addr, unsigned flags);
extern void __uv_udp_send_cb(uv_udp_send_t* req, int status);
extern void __uv_timer_cb(uv_timer_t* timer, int status);
extern void __uv_idle_cb(uv_idle_t* handle, int status);
extern void __uv_close_cb(uv_handle_t* handle);
extern void __uv_shutdown_cb(uv_shutdown_t* req, int status);

static uv_buf_t _uv_alloc_cb(uv_handle_t* handle, size_t suggested_size) {
    char* buf;
    buf = (char*)malloc(suggested_size);
    return uv_buf_init(buf, suggested_size);
}

static int _uv_udp_send(uv_udp_send_t* req, uv_udp_t* handle, uv_buf_t bufs[], int bufcnt, struct sockaddr_in addr) {
	return uv_udp_send(req, handle, bufs, bufcnt, addr, __uv_udp_send_cb);
}

static int _uv_udp_send6(uv_udp_send_t* req, uv_udp_t* handle, uv_buf_t bufs[], int bufcnt, struct sockaddr_in6 addr) {
	return uv_udp_send6(req, handle, bufs, bufcnt, addr, __uv_udp_send_cb);
}

static int _uv_udp_recv_start(uv_udp_t* udp) {
	return uv_udp_recv_start(udp, _uv_alloc_cb, __uv_udp_recv_cb);
}

static int _uv_tcp_connect(uv_connect_t* req, uv_tcp_t* handle, struct sockaddr_in address) {
	return uv_tcp_connect(req, handle, address, __uv_connect_cb);
}

static int _uv_tcp_connect6(uv_connect_t* req, uv_tcp_t* handle, struct sockaddr_in6 address) {
	return uv_tcp_connect6(req, handle, address, __uv_connect_cb);
}

static void _uv_pipe_connect(uv_connect_t* req, uv_pipe_t* handle, const char* name) {
	uv_pipe_connect(req, handle, name, __uv_connect_cb);
}

static int _uv_listen(uv_stream_t* stream, int backlog) {
	return uv_listen(stream, backlog, __uv_connection_cb);
}

static int _uv_read_start(uv_stream_t* stream) {
	return uv_read_start(stream, _uv_alloc_cb, __uv_read_cb);
}

static int _uv_write(uv_write_t* req, uv_stream_t* handle, uv_buf_t bufs[], int bufcnt) {
	return uv_write(req, handle, bufs, bufcnt, __uv_write_cb);
}

static void _uv_close(uv_handle_t* handle) {
	uv_close(handle, __uv_close_cb);
}

static int _uv_shutdown(uv_shutdown_t* req, uv_stream_t* handle) {
	return uv_shutdown(req, handle, __uv_shutdown_cb);
}

static int _uv_timer_start(uv_timer_t* timer, int64_t timeout, int64_t repeat) {
	return uv_timer_start(timer, __uv_timer_cb, timeout, repeat);
}

static int _uv_idle_start(uv_idle_t* idle) {
	return uv_idle_start(idle, __uv_idle_cb);
}

//static uv_stream_t* _uv_to_stream(void* any) {
//	return (uv_stream_t*) any;
//}
//
//static uv_tcp_t* _uv_to_tcp(void* any) {
//	return (uv_tcp_t*) any;
//}
//
//static uv_pipe_t* _uv_to_pipe(void* any) {
//	return (uv_pipe_t*) any;
//}
//
//static uv_udp_t* _uv_to_udp(void* any) {
//	return (uv_udp_t*) any;
//}
//
//static uv_timer_t* _uv_to_timer(void* any) {
//	return (uv_timer_t*) any;
//}
//
//static uv_idle_t* _uv_to_idle(void* any) {
//	return (uv_idle_t*) any;
//}
//
//static struct sockaddr* _uv_to_sa(void* any) {
//	return (struct sockaddr*) any;
//}
//
//static struct sockaddr_in _uv_to_sai4(void* any) {
//	return *(struct sockaddr_in*) any;
//}
//
//static struct sockaddr_in6 _uv_to_sai6(void* any) {
//	return *(struct sockaddr_in6*) any;
//}
//
//static uv_handle_t* _uv_to_handle(void* any) {
//	return (uv_handle_t*) any;
//}

#define UV_SIZEOF_SOCKADDR_IN ((int)sizeof(struct sockaddr_in))

#cgo darwin LDFLAGS: -luv
#cgo linux LDFLAGS: -ldl -luv -lpthread -lrt -lm
#cgo windows LDFLAGS: -luv.dll -lws2_32
*/
import "C"
import "unsafe"

type Request struct {
	r      *C.uv_req_t
	Handle *Handle
}

type Handle struct {
	h    *C.uv_handle_t
	Data interface{}
}

type callback_info struct {
	connection_cb func(*Handle, int)
	connect_cb    func(*Request, int)
	read_cb       func(*Handle, []byte)
	udp_recv_cb   func(*Handle, []byte, SockaddrIn, uint)
	write_cb      func(*Request, int)
	udp_send_cb   func(*Request, int)
	close_cb      func(*Handle)
	shutdown_cb   func(*Request, int)
	timer_cb      func(*Handle, int)
	idle_cb       func(*Handle, int)
	data          interface{}
}

func uv_tcp_bind(tcp *C.uv_tcp_t, sa4 C.struct_sockaddr_in) int {
	return int(C.uv_tcp_bind(tcp, sa4))
}

func uv_tcp_bind6(tcp *C.uv_tcp_t, sa6 C.struct_sockaddr_in6) int {
	return int(C.uv_tcp_bind6(tcp, sa6))
}

func uv_tcp_connect(tcp *C.uv_tcp_t, sa4 C.struct_sockaddr_in) int {
	var req C.uv_connect_t
	return int(C._uv_tcp_connect(&req, tcp, sa4))
}

func uv_tcp_connect6(tcp *C.uv_tcp_t, sa6 C.struct_sockaddr_in6) int {
	var req C.uv_connect_t
	return int(C._uv_tcp_connect6(&req, tcp, sa6))
}

func uv_pipe_connect(pipe *C.uv_pipe_t, name string) {
	pname := C.CString(name)
	defer C.free(unsafe.Pointer(pname))
	var req C.uv_connect_t
	C._uv_pipe_connect(&req, pipe, pname)
}

func uv_pipe_bind(pipe *C.uv_pipe_t, name string) int {
	pname := C.CString(name)
	defer C.free(unsafe.Pointer(pname))
	return int(C.uv_pipe_bind(pipe, pname))
}

func uv_close(handle *C.uv_handle_t) {
	C._uv_close(handle)
}

func uv_is_active(handle *C.uv_handle_t) bool {
	if C.uv_is_active(handle) != 0 {
		return true
	}
	return false
}

func uv_listen(stream *C.uv_stream_t, backlog int) int {
	return int(C._uv_listen(stream, C.int(backlog)))
}

func uv_accept(stream *C.uv_stream_t, client *C.uv_stream_t) int {
	return int(C.uv_accept(stream, client))
}

func uv_shutdown(stream *C.uv_stream_t) int {
	var req C.uv_shutdown_t
	return int(C._uv_shutdown(&req, stream))
}

func uv_write(stream *C.uv_stream_t, buf *C.uv_buf_t, bufcnt int) int {
	var req C.uv_write_t
	return int(C._uv_write(&req, stream, buf, C.int(bufcnt)))
}

func uv_udp_send(udp *C.uv_udp_t, buf *C.uv_buf_t, bufcnt int, sa4 C.struct_sockaddr_in) int {
	var req C.uv_udp_send_t
	return int(C._uv_udp_send(&req, udp, buf, C.int(bufcnt), sa4))
}

func uv_udp_send6(udp *C.uv_udp_t, buf *C.uv_buf_t, bufcnt int, sa6 C.struct_sockaddr_in6) int {
	var req C.uv_udp_send_t
	return int(C._uv_udp_send6(&req, udp, buf, C.int(bufcnt), sa6))
}

func uv_read_start(stream *C.uv_stream_t) int {
	return int(C._uv_read_start(stream))
}

func uv_read_stop(stream *C.uv_stream_t) int {
	return int(C.uv_read_stop(stream))
}

func uv_udp_recv_start(udp *C.uv_udp_t) int {
	return int(C._uv_udp_recv_start(udp))
}

func uv_udp_recv_stop(udp *C.uv_udp_t) int {
	return int(C.uv_udp_recv_stop(udp))
}

func uv_buf_init(b []byte) C.uv_buf_t {
	return C.uv_buf_init((*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)))
}

func uv_tcp_nodelay(tcp *C.uv_tcp_t, enable bool) int {
	var v C.int
	if enable {
		v = 1
	}
	return int(C.uv_tcp_nodelay(tcp, v))
}

func uv_tcp_keepalive(tcp *C.uv_tcp_t, enable bool, delay uint) int {
	var v C.int
	if enable {
		v = 1
	}
	return int(C.uv_tcp_keepalive(tcp, v, C.uint(delay)))
}

func uv_tcp_simultaneous_accepts(tcp *C.uv_tcp_t, enable bool) int {
	var v C.int
	if enable {
		v = 1
	}
	return int(C.uv_tcp_simultaneous_accepts(tcp, v))
}

func uv_tcp_getsockname(tcp *C.uv_tcp_t, sa *C.struct_sockaddr) int {
	l := C.UV_SIZEOF_SOCKADDR_IN
	return int(C.uv_tcp_getsockname(tcp, sa, (*C.int)(unsafe.Pointer(&l))))
}

func uv_tcp_getpeername(tcp *C.uv_tcp_t, sa *C.struct_sockaddr) int {
	l := C.UV_SIZEOF_SOCKADDR_IN
	return int(C.uv_tcp_getpeername(tcp, sa, (*C.int)(unsafe.Pointer(&l))))
}

func uv_udp_getsockname(udp *C.uv_udp_t, sa *C.struct_sockaddr) int {
	l := C.UV_SIZEOF_SOCKADDR_IN
	return int(C.uv_udp_getsockname(udp, sa, (*C.int)(unsafe.Pointer(&l))))
}

func uv_timer_start(timer *C.uv_timer_t, timeout int64, repeat int64) int {
	return int(C._uv_timer_start(timer, C.int64_t(timeout), C.int64_t(repeat)))
}

func uv_idle_start(idle *C.uv_idle_t) int {
	return int(C._uv_idle_start(idle))
}

//export __uv_connect_cb
func __uv_connect_cb(p unsafe.Pointer, status int) {
	c := (*C.uv_connect_t)(p)
	cbi := (*callback_info)(c.handle.data)
	if cbi.connect_cb != nil {
		cbi.connect_cb(&Request{
			(*C.uv_req_t)(unsafe.Pointer(c)),
			&Handle{
				(*C.uv_handle_t)(unsafe.Pointer(c.handle)),
				cbi.data}}, status)
	}
}

//export __uv_connection_cb
func __uv_connection_cb(p unsafe.Pointer, status int) {
	s := (*C.uv_stream_t)(p)
	cbi := (*callback_info)(s.data)
	if cbi.connection_cb != nil {
		cbi.connection_cb(&Handle{(*C.uv_handle_t)(unsafe.Pointer(s)), cbi.data}, status)
	}
}

//export __uv_read_cb
func __uv_read_cb(p unsafe.Pointer, nread int, buf unsafe.Pointer) {
	s := (*C.uv_stream_t)(p)
	cbi := (*callback_info)(s.data)
	if cbi.read_cb != nil {
		cbi.read_cb(&Handle{(*C.uv_handle_t)(unsafe.Pointer(s)), cbi.data}, (*[1 << 30]byte)(unsafe.Pointer(buf))[0:nread])
	}
}

//export __uv_write_cb
func __uv_write_cb(p unsafe.Pointer, status int) {
	w := (*C.uv_write_t)(p)
	cbi := (*callback_info)(w.handle.data)
	if cbi.write_cb != nil {
		cbi.write_cb(&Request{
			(*C.uv_req_t)(unsafe.Pointer(w)),
			&Handle{
				(*C.uv_handle_t)(unsafe.Pointer(w.handle)),
				cbi.data}}, status)
	}
}

//export __uv_close_cb
func __uv_close_cb(p unsafe.Pointer) {
	h := (*C.uv_handle_t)(p)
	cbi := (*callback_info)(h.data)
	if cbi.close_cb != nil {
		cbi.close_cb(&Handle{h, cbi.data})
	}
}

//export __uv_shutdown_cb
func __uv_shutdown_cb(p unsafe.Pointer, status int) {
	s := (*C.uv_shutdown_t)(p)
	cbi := (*callback_info)(s.handle.data)
	if cbi.shutdown_cb != nil {
		cbi.shutdown_cb(&Request{
			(*C.uv_req_t)(unsafe.Pointer(s)),
			&Handle{
				(*C.uv_handle_t)(unsafe.Pointer(s.handle)),
				cbi.data}}, status)
	}
}

//export __uv_udp_recv_cb
func __uv_udp_recv_cb(p unsafe.Pointer, nread int, buf unsafe.Pointer, sa unsafe.Pointer, flags uint) {
	u := (*C.uv_udp_t)(p)
	cbi := (*callback_info)(u.data)
	if cbi.udp_recv_cb != nil {
		psa := &SockaddrIn4{*(*C.struct_sockaddr_in)(sa)}
		cbi.udp_recv_cb(&Handle{
			(*C.uv_handle_t)(unsafe.Pointer(u)), cbi.data}, (*[1 << 30]byte)(unsafe.Pointer(buf))[0:nread], psa, flags)
	}
}

//export __uv_udp_send_cb
func __uv_udp_send_cb(p unsafe.Pointer, status int) {
	us := (*C.uv_udp_send_t)(p)
	cbi := (*callback_info)(us.handle.data)
	if cbi.udp_send_cb != nil {
		cbi.udp_send_cb(&Request{
			(*C.uv_req_t)(unsafe.Pointer(us)),
			&Handle{
				(*C.uv_handle_t)(unsafe.Pointer(us.handle)),
				cbi.data}}, status)
	}
}

//export __uv_timer_cb
func __uv_timer_cb(p unsafe.Pointer, status int) {
	t := (*C.uv_timer_t)(p)
	cbi := (*callback_info)(t.data)
	if cbi.timer_cb != nil {
		cbi.timer_cb(&Handle{
			(*C.uv_handle_t)(unsafe.Pointer(t)), cbi.data}, status)
	}
}

//export __uv_idle_cb
func __uv_idle_cb(p unsafe.Pointer, status int) {
	i := (*C.uv_idle_t)(p)
	cbi := (*callback_info)(i.data)
	if cbi.idle_cb != nil {
		cbi.idle_cb(&Handle{
			(*C.uv_handle_t)(unsafe.Pointer(i)), cbi.data}, status)
	}
}
