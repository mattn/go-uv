package uv

/*
#include <uv/uv.h>

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

static int _uv_listen(uv_stream_t* stream, int backlog) {
	return uv_listen(stream, backlog, _uv_connection_cb);
}

static void _uv_pipe_connect(uv_connect_t* req, uv_pipe_t* handle, const char* name) {
	return uv_pipe_connect(req, handle, name, _uv_connect_cb);
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

static uv_stream_t* _uv_pipe_to_stream(uv_pipe_t* pipe) {
	return (uv_stream_t*) pipe;
}

static uv_handle_t* _uv_pipe_to_handle(uv_pipe_t* pipe) {
	return (uv_handle_t*) pipe;
}
*/
import "C"
import "errors"
import "unsafe"

type Pipe struct {
	p *C.uv_pipe_t
}

type pipe_callback_info struct {
	connection_cb func(int)
	connect_cb func(int)
	pipe_recv_cb       func([]byte, SockaddrIn, uint)
	write_cb      func(int)
	close_cb      func()
	shutdown_cb    func(int)
}

func PipeInit(ipc int) (pipe *Pipe, err error) {
	var p C.uv_pipe_t

	r := C.uv_pipe_init(C.uv_default_loop(), &p, C.int(ipc))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return nil, errors.New(C.GoString(C.uv_strerror(e)))
	}
	p.data = unsafe.Pointer(&pipe_callback_info{})
	return &Pipe{&p}, nil
}

func (pipe *Pipe) Bind(name string) (err error) {
	pname := C.CString(name)
	defer C.free(unsafe.Pointer(pname))
	r := C.uv_pipe_bind(pipe.p, pname)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (pipe *Pipe) Listen(backlog int, cb func(int)) (err error) {
	cbi := (*pipe_callback_info)(pipe.p.data)
	cbi.connection_cb = cb
	r := C._uv_listen(C._uv_pipe_to_stream(pipe.p), C.int(backlog))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (pipe *Pipe) Connect(name string, cb func(int)) (err error) {
	pname := C.CString(name)
	defer C.free(unsafe.Pointer(pname))
	cbi := (*pipe_callback_info)(pipe.p.data)
	cbi.connect_cb = cb
	var req C.uv_connect_t
	C._uv_pipe_connect(&req, pipe.p, pname)
	// TODO: error
	return nil
}

func (pipe *Pipe) Accept() (client *Pipe, err error) {
	c, err := PipeInit(1)
	if err != nil {
		return nil, err
	}
	r := C.uv_accept(C._uv_pipe_to_stream(pipe.p), C._uv_pipe_to_stream(c.p))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return nil, errors.New(C.GoString(C.uv_strerror(e)))
	}
	return &Pipe{c.p}, nil
}

func (pipe *Pipe) RecvStart(cb func([]byte, SockaddrIn, uint)) (err error) {
	cbi := (*pipe_callback_info)(pipe.p.data)
	cbi.pipe_recv_cb = cb
	r := C._uv_read_start(C._uv_pipe_to_stream(pipe.p))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (pipe *Pipe) RecvStop() (err error) {
	r := C.uv_read_stop(C._uv_pipe_to_stream(pipe.p))
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (pipe *Pipe) Write(b []byte, cb func(int)) (err error) {
	cbi := (*pipe_callback_info)(pipe.p.data)
	cbi.write_cb = cb
	var req C.uv_write_t
	buf := C.uv_buf_init((*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)))
	r := C._uv_write(&req, C._uv_pipe_to_stream(pipe.p), &buf, 1)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (pipe *Pipe) Shutdown(cb func(int)) {
	cbi := (*pipe_callback_info)(pipe.p.data)
	cbi.shutdown_cb = cb
	var req C.uv_shutdown_t
	C._uv_shutdown(&req, C._uv_pipe_to_stream(pipe.p))
}

func (pipe *Pipe) Close(cb func()) {
	cbi := (*pipe_callback_info)(pipe.p.data)
	cbi.close_cb = cb
	C._uv_close(C._uv_pipe_to_handle(pipe.p))
}

func (pipe *Pipe) IsActive() bool {
	if C.uv_is_active(C._uv_pipe_to_handle(pipe.p)) != 0 {
		return true
	}
	return false
}
