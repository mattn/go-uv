package uv

/*
#include <stdio.h>
#include <stdlib.h>
#include <uv/uv.h>

//========== TCP ==========
extern void __uv_connect_cb(void* p, int status);
static void _uv_connect_cb(uv_connect_t* req, int status) {
	__uv_connect_cb(req->handle->data, status);
}

extern void __uv_connection_cb(void* p, int status);
static void _uv_connection_cb(uv_stream_t* server, int status) {
	__uv_connection_cb(server->data, status);
}

extern void __uv_read_cb(void* p, int nread, void* buf);
static void _uv_read_cb(uv_stream_t* stream, ssize_t nread, uv_buf_t buf) {
	__uv_read_cb(stream->data, nread, buf.base);
}

extern void __uv_write_cb(void* p, int status);
static void _uv_write_cb(uv_write_t* req, int status) {
	__uv_write_cb(req->handle->data, status);
}

static uv_buf_t _uv_alloc_cb(uv_handle_t* handle, size_t suggested_size) {
    char* buf;
    buf = (char*)malloc(suggested_size);
    return uv_buf_init(buf, suggested_size);
}

extern void __uv_close_cb(void* p);
static void _uv_close_cb(uv_handle_t* handle) {
	__uv_close_cb(handle->data);
}

static int _uv_listen(uv_stream_t* stream, int backlog) {
	return uv_listen(stream, backlog, _uv_connection_cb);
}

static int _uv_tcp_connect(uv_connect_t* req, uv_tcp_t* handle, struct sockaddr_in address) {
	return uv_tcp_connect(req, handle, address, _uv_connect_cb);
}

static int _uv_tcp_connect6(uv_connect_t* req, uv_tcp_t* handle, struct sockaddr_in6 address) {
	return uv_tcp_connect6(req, handle, address, _uv_connect_cb);
}

static int _uv_read_start(uv_stream_t* stream) {
	return uv_read_start(stream, _uv_alloc_cb, _uv_read_cb);
}

static int _uv_write(uv_write_t* req, uv_stream_t* handle, uv_buf_t bufs[], int bufcnt) {
	return uv_write(req, handle, bufs, bufcnt, _uv_write_cb);
}

static void _uv_close(uv_handle_t* handle) {
	uv_close(handle, _uv_close_cb);
}

//========== TIMER ==========
extern void __uv_timer_cb(void* p, int status);
static void _uv_timer_cb(uv_timer_t* handle, int status) {
	__uv_timer_cb(handle->data, status);
}

static int _uv_timer_start(uv_timer_t* handle, int64_t timeout, int64_t repeat) {
	return uv_timer_start(handle, _uv_timer_cb, timeout, repeat);
}

//========== UTIL ==========
static uv_stream_t* _uv_tcp_to_stream(uv_tcp_t* tcp) {
	return (uv_stream_t*) tcp;
}

static uv_handle_t* _uv_tcp_to_handle(uv_tcp_t* tcp) {
	return (uv_handle_t*) tcp;
}

#cgo darwin LDFLAGS: -luv
#cgo linux LDFLAGS: -luv
#cgo windows LDFLAGS: -luv.dll -lws2_32
*/
import "C"
import "fmt"
import "errors"
import "unsafe"

type Tcp struct {
	t *C.uv_tcp_t
}

type tcp_callback_info struct {
	connection_cb func(int)
	connect_cb    func(int)
	read_cb       func([]byte)
	write_cb      func(int)
	close_cb      func()
}

func TcpInit() (tcp *Tcp, err error) {
	var t C.uv_tcp_t

	r := C.uv_tcp_init(C.uv_default_loop(), &t)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return nil, errors.New(C.GoString(C.uv_strerror(e)))
	}
	t.data = unsafe.Pointer(&tcp_callback_info{})
	return &Tcp{&t}, nil
}

func (tcp *Tcp) Bind(host string, port uint16) (err error) {
	phost := C.CString(host)
	defer C.free(unsafe.Pointer(phost))
	r := C.uv_tcp_bind(tcp.t, C.uv_ip4_addr(phost, C.int(port)))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (tcp *Tcp) Bind6(host string, port uint16) (err error) {
	phost := C.CString(host)
	defer C.free(unsafe.Pointer(phost))
	r := C.uv_tcp_bind6(tcp.t, C.uv_ip6_addr(phost, C.int(port)))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (tcp *Tcp) Nodelay(enable bool) (err error) {
	var v C.int
	if enable {
		v = 1
	}
	r := C.uv_tcp_nodelay(tcp.t, v)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (tcp *Tcp) Keepalive(enable bool, delay uint) (err error) {
	var v C.int
	if enable {
		v = 1
	}
	r := C.uv_tcp_keepalive(tcp.t, v, C.uint(delay))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
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
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (tcp *Tcp) Connect(host string, port uint16, cb func(int)) (err error) {
	cbi := (*tcp_callback_info)(tcp.t.data)
	cbi.connect_cb = cb
	phost := C.CString(host)
	defer C.free(unsafe.Pointer(phost))
	var req C.uv_connect_t
	r := C._uv_tcp_connect(&req, tcp.t, C.uv_ip4_addr(phost, C.int(port)))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (tcp *Tcp) Connect6(host string, port uint16, cb func(int)) (err error) {
	cbi := (*tcp_callback_info)(tcp.t.data)
	cbi.connect_cb = cb
	phost := C.CString(host)
	defer C.free(unsafe.Pointer(phost))
	var req C.uv_connect_t
	r := C._uv_tcp_connect6(&req, tcp.t, C.uv_ip6_addr(phost, C.int(port)))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (tcp *Tcp) Listen(backlog int, cb func(int)) (err error) {
	cbi := (*tcp_callback_info)(tcp.t.data)
	cbi.connection_cb = cb
	r := C._uv_listen(C._uv_tcp_to_stream(tcp.t), C.int(backlog))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (tcp *Tcp) Accept() (client *Tcp, err error) {
	c, err := TcpInit()
	if err != nil {
		return nil, err
	}
	r := C.uv_accept(C._uv_tcp_to_stream(tcp.t), C._uv_tcp_to_stream(c.t))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return nil, errors.New(C.GoString(C.uv_strerror(e)))
	}
	return &Tcp{c.t}, nil
}

func (tcp *Tcp) ReadStart(cb func([]byte)) (err error) {
	cbi := (*tcp_callback_info)(tcp.t.data)
	cbi.read_cb = cb
	r := C._uv_read_start(C._uv_tcp_to_stream(tcp.t))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (tcp *Tcp) ReadStop() (err error) {
	r := C.uv_read_stop(C._uv_tcp_to_stream(tcp.t))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (tcp *Tcp) Write(b []byte, cb func(int)) (err error) {
	cbi := (*tcp_callback_info)(tcp.t.data)
	cbi.write_cb = cb
	var req C.uv_write_t
	buf := C.uv_buf_init((*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)))
	r := C._uv_write(&req, C._uv_tcp_to_stream(tcp.t), &buf, 1)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (tcp *Tcp) Close(cb func()) {
	cbi := (*tcp_callback_info)(tcp.t.data)
	cbi.close_cb = cb
	C._uv_close(C._uv_tcp_to_handle(tcp.t))
}

func (tcp *Tcp) IsActive() bool {
	if C.uv_is_active(C._uv_tcp_to_handle(tcp.t)) != 0 {
		return true
	}
	return false
}

//export __uv_connect_cb
func __uv_connect_cb(p unsafe.Pointer, status int) {
	cbi := (*tcp_callback_info)(p)
	if cbi.connect_cb != nil {
		cbi.connect_cb(status)
	}
}

//export __uv_connection_cb
func __uv_connection_cb(p unsafe.Pointer, status int) {
	cbi := (*tcp_callback_info)(p)
	if cbi.connection_cb != nil {
		cbi.connection_cb(status)
	}
}

//export __uv_read_cb
func __uv_read_cb(p unsafe.Pointer, nread int, buf unsafe.Pointer) {
	cbi := (*tcp_callback_info)(p)
	if cbi.read_cb != nil {
		cbi.read_cb((*[1 << 30]byte)(unsafe.Pointer(buf))[0:nread])
	}
}

//export __uv_write_cb
func __uv_write_cb(p unsafe.Pointer, status int) {
	cbi := (*tcp_callback_info)(p)
	if cbi.write_cb != nil {
		cbi.write_cb(status)
	}
}

//export __uv_close_cb
func __uv_close_cb(p unsafe.Pointer) {
	cbi := (*tcp_callback_info)(p)
	if cbi.close_cb != nil {
		cbi.close_cb()
	}
}

type Loop struct {
	l *C.uv_loop_t
}

func DefaultLoop() *Loop {
	return &Loop{C.uv_default_loop()}
}

func LoopNew() *Loop {
	return &Loop{C.uv_loop_new()}
}

func (loop *Loop) Delete() {
	C.uv_loop_delete(loop.l)
}

func (loop *Loop) Run() {
	C.uv_run(loop.l)
}

func (loop *Loop) Ref() {
	C.uv_ref(loop.l)
}

func (loop *Loop) Unref() {
	C.uv_unref(loop.l)
}

func (loop *Loop) UpdateTime() {
	C.uv_update_time(loop.l)
}

func (loop *Loop) Now() int64 {
	return int64(C.uv_now(loop.l))
}

func Version() string {
    return fmt.Sprintf("%d.%d", C.UV_VERSION_MAJOR, C.UV_VERSION_MINOR)
}

type Error struct {
	e C.uv_err_t
}

func LastError() *Error {
	return &Error{C.uv_last_error(C.uv_default_loop())}
}

func (err *Error) String() string {
	return C.GoString(C.uv_strerror(err.e))
}

func (err *Error) Name() string {
	return C.GoString(C.uv_err_name(err.e))
}

type Timer struct {
	t *C.uv_timer_t
}

type timer_callback_info struct {
	timer_cb      func(int)
}

func TimerInit() (timer *Timer, err error) {
	var t C.uv_timer_t

	r := C.uv_timer_init(C.uv_default_loop(), &t)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return nil, errors.New(C.GoString(C.uv_strerror(e)))
	}
	t.data = unsafe.Pointer(&timer_callback_info{})
	return &Timer{&t}, nil
}

func (timer *Timer) Start(timeout int64, repeat int64, cb func(int)) (err error) {
	cbi := (*timer_callback_info)(timer.t.data)
	cbi.timer_cb = cb
	r := C._uv_timer_start(timer.t, C.int64_t(timeout), C.int64_t(repeat))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (timer *Timer) Stop() (err error) {
	r := C.uv_timer_stop(timer.t)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (timer *Timer) Again() (err error) {
	r := C.uv_timer_again(timer.t)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (timer *Timer) SetRepeat(repeat int64) {
	C.uv_timer_set_repeat(timer.t, C.int64_t(repeat))
}

func (timer *Timer) GetRepeat() int64 {
	return int64(C.uv_timer_get_repeat(timer.t))
}

//export __uv_timer_cb
func __uv_timer_cb(p unsafe.Pointer, status int) {
	cbi := (*timer_callback_info)(p)
	if cbi.timer_cb != nil {
		cbi.timer_cb(status)
	}
}

