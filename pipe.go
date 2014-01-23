package uv

/*
#include <uv.h>
*/
import "C"
import "unsafe"

type Pipe struct {
	p *C.uv_pipe_t
	l *C.uv_loop_t
	Handle
}

func PipeInit(loop *Loop, ipc int) (pipe *Pipe, err error) {
	var p C.uv_pipe_t

	if loop == nil {
		loop = DefaultLoop()
	}
	r := C.uv_pipe_init(loop.l, &p, C.int(ipc))
	if r != 0 {
		return nil, &Error{int(r)}
	}
	p.data = unsafe.Pointer(&callback_info{})
	return &Pipe{&p, loop.l, Handle{(*C.uv_handle_t)(unsafe.Pointer(&p)), p.data}}, nil
}

func (pipe *Pipe) GetLoop() *Loop {
	return &Loop{pipe.l}
}

func (pipe *Pipe) Open(name string) (err error) {
	r := uv_pipe_bind(pipe.p, name)
	if r != 0 {
		return &Error{r}
	}
	return nil
}

func (pipe *Pipe) Bind(name string) (err error) {
	r := uv_pipe_bind(pipe.p, name)
	if r != 0 {
		return &Error{r}
	}
	return nil
}

func (pipe *Pipe) Listen(backlog int, cb func(*Handle, int)) (err error) {
	cbi := (*callback_info)(pipe.p.data)
	cbi.connection_cb = cb
	r := uv_listen((*C.uv_stream_t)(unsafe.Pointer(pipe.p)), backlog)
	if r != 0 {
		return &Error{r}
	}
	return nil
}

func (pipe *Pipe) Connect(name string, cb func(*Request, int)) (err error) {
	cbi := (*callback_info)(pipe.p.data)
	cbi.connect_cb = cb
	uv_pipe_connect(pipe.p, name)
	// TODO: error
	return nil
}

func (pipe *Pipe) Accept() (client *Pipe, err error) {
	c, err := PipeInit(pipe.GetLoop(), 0)
	if err != nil {
		return nil, err
	}
	r := uv_accept((*C.uv_stream_t)(unsafe.Pointer(pipe.p)), (*C.uv_stream_t)(unsafe.Pointer(c.p)))
	if r != 0 {
		return nil, &Error{r}
	}
	return &Pipe{c.p, pipe.l, Handle{(*C.uv_handle_t)(unsafe.Pointer(c.p)), c.p.data}}, nil
}

func (pipe *Pipe) ReadStart(cb func(*Handle, []byte)) (err error) {
	cbi := (*callback_info)(pipe.p.data)
	cbi.read_cb = cb
	r := uv_read_start((*C.uv_stream_t)(unsafe.Pointer(pipe.p)))
	if r != 0 {
		return &Error{r}
	}
	return nil
}

func (pipe *Pipe) ReadStop() (err error) {
	r := uv_read_stop((*C.uv_stream_t)(unsafe.Pointer(pipe.p)))
	if r != 0 {
		return &Error{r}
	}
	return nil
}

func (pipe *Pipe) Write(b []byte, cb func(*Request, int)) (err error) {
	cbi := (*callback_info)(pipe.p.data)
	cbi.write_cb = cb
	buf := C.uv_buf_init((*C.char)(unsafe.Pointer(&b[0])), C.uint(len(b)))
	r := uv_write((*C.uv_stream_t)(unsafe.Pointer(pipe.p)), &buf, 1)
	if r != 0 {
		return &Error{r}
	}
	return nil
}

func (pipe *Pipe) Shutdown(cb func(*Request, int)) {
	cbi := (*callback_info)(pipe.p.data)
	cbi.shutdown_cb = cb
	uv_shutdown((*C.uv_stream_t)(unsafe.Pointer(pipe.p)))
}
