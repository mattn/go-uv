package uv

/*
#include <uv/uv.h>
*/
import "C"
import "unsafe"

type Idle struct {
	i *C.uv_idle_t
	l *C.uv_loop_t
	Handle
}

func IdleInit(loop *Loop) (idle *Idle, err error) {
	var i C.uv_idle_t

	if loop == nil {
		loop = DefaultLoop()
	}
	r := C.uv_idle_init(loop.l, &i)
	if r != 0 {
		return nil, idle.GetLoop().LastError().Error()
	}
	i.data = unsafe.Pointer(&callback_info{})
	return &Idle{&i, loop.l, Handle{(*C.uv_handle_t)(unsafe.Pointer(&i)), i.data}}, nil
}

func (idle *Idle) GetLoop() *Loop {
	return &Loop{idle.l}
}

func (idle *Idle) Start(cb func(*Handle, int)) (err error) {
	cbi := (*callback_info)(idle.i.data)
	cbi.idle_cb = cb
	r := uv_idle_start(idle.i)
	if r != 0 {
		return idle.GetLoop().LastError().Error()
	}
	return nil
}

func (idle *Idle) Stop() (err error) {
	r := C.uv_idle_stop(idle.i)
	if r != 0 {
		return idle.GetLoop().LastError().Error()
	}
	return nil
}
