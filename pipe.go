package uv

/*
#include <uv/uv.h>
*/
import "C"
import "unsafe"

type Pipe struct {
	p *C.uv_pipe_t
	l *C.uv_loop_t
}

func PipeInit(loop *Loop, ipc int) (pipe *Pipe, err error) {
	var p C.uv_pipe_t

	if loop == nil {
		loop = DefaultLoop()
	}
	r := C.uv_pipe_init(loop.l, &p, C.int(ipc))
	if r != 0 {
		return nil, pipe.GetLoop().LastError().Error()
	}
	p.data = unsafe.Pointer(&callback_info{})
	return &Pipe{&p, loop.l}, nil
}

func (pipe *Pipe) GetLoop() *Loop {
	return &Loop{pipe.l}
}

func (pipe *Pipe) Bind(name string) (err error) {
	r := uv_pipe_bind(pipe.p, name)
	if r != 0 {
		return pipe.GetLoop().LastError().Error()
	}
	return nil
}

func (pipe *Pipe) Listen(backlog int, cb func(*Handle, int)) (err error) {
	cbi := (*callback_info)(pipe.p.data)
	cbi.connection_cb = cb
	r := uv_listen((*C.uv_stream_t)(unsafe.Pointer(pipe.p)), backlog)
	if r != 0 {
		return pipe.GetLoop().LastError().Error()
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
	c, err := PipeInit(pipe.GetLoop(), 1)
	if err != nil {
		return nil, err
	}
	r := uv_accept((*C.uv_stream_t)(unsafe.Pointer(pipe.p)), (*C.uv_stream_t)(unsafe.Pointer(c.p)))
	if r != 0 {
		return nil, pipe.GetLoop().LastError().Error()
	}
	return &Pipe{c.p, pipe.l}, nil
}

func (pipe *Pipe) RecvStart(cb func(*Handle, []byte)) (err error) {
	cbi := (*callback_info)(pipe.p.data)
	cbi.read_cb = cb
	r := uv_read_start((*C.uv_stream_t)(unsafe.Pointer(pipe.p)))
	if r != 0 {
		return pipe.GetLoop().LastError().Error()
	}
	return nil
}

func (pipe *Pipe) RecvStop() (err error) {
	r := uv_read_stop((*C.uv_stream_t)(unsafe.Pointer(pipe.p)))
	if r != 0 {
		return pipe.GetLoop().LastError().Error()
	}
	return nil
}

func (pipe *Pipe) Write(b []byte, cb func(*Request, int)) (err error) {
	cbi := (*callback_info)(pipe.p.data)
	cbi.write_cb = cb
	buf := C.uv_buf_init((*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)))
	r := uv_write((*C.uv_stream_t)(unsafe.Pointer(pipe.p)), &buf, 1)
	if r != 0 {
		return pipe.GetLoop().LastError().Error()
	}
	return nil
}

func (pipe *Pipe) Shutdown(cb func(*Request, int)) {
	cbi := (*callback_info)(pipe.p.data)
	cbi.shutdown_cb = cb
	uv_shutdown((*C.uv_stream_t)(unsafe.Pointer(pipe.p)))
}

func (pipe *Pipe) Close(cb func(*Handle)) {
	cbi := (*callback_info)(pipe.p.data)
	cbi.close_cb = cb
	uv_close((*C.uv_handle_t)(unsafe.Pointer(pipe.p)))
}

func (pipe *Pipe) IsActive() bool {
	return uv_is_active((*C.uv_handle_t)(unsafe.Pointer(pipe.p)))
}
