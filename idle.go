package uv

/*
#include <uv/uv.h>
*/
import "C"
import "unsafe"

type Idle struct {
	i *C.uv_idle_t
	l *C.uv_loop_t
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
	return &Idle{&i, loop.l}, nil
}

func (idle *Idle) GetLoop() *Loop {
	return &Loop{idle.l}
}

func (idle *Idle) Start(cb func(int)) (err error) {
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

func (idle *Idle) Close(cb func()) {
	cbi := (*callback_info)(idle.i.data)
	cbi.close_cb = cb
	uv_close((*C.uv_handle_t)(unsafe.Pointer(idle.i)))
}

func (idle *Idle) IsActive() bool {
	return uv_is_active((*C.uv_handle_t)(unsafe.Pointer(idle.i)))
}
