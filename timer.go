package uv

/*
#include <uv/uv.h>
*/
import "C"
import "unsafe"

type Timer struct {
	t *C.uv_timer_t
	l *C.uv_loop_t
	Handle
}

func TimerInit(loop *Loop) (timer *Timer, err error) {
	var t C.uv_timer_t

	if loop == nil {
		loop = DefaultLoop()
	}
	r := C.uv_timer_init(loop.l, &t)
	if r != 0 {
		return nil, timer.GetLoop().LastError().Error()
	}
	t.data = unsafe.Pointer(&callback_info{})
	return &Timer{&t, loop.l, Handle{(*C.uv_handle_t)(unsafe.Pointer(&t)), t.data}}, nil
}

func (timer *Timer) GetLoop() *Loop {
	return &Loop{timer.l}
}

func (timer *Timer) Start(cb func(*Handle, int), timeout int64, repeat int64) (err error) {
	cbi := (*callback_info)(timer.t.data)
	cbi.timer_cb = cb
	r := uv_timer_start(timer.t, timeout, repeat)
	if r != 0 {
		return timer.GetLoop().LastError().Error()
	}
	return nil
}

func (timer *Timer) Stop() (err error) {
	r := C.uv_timer_stop(timer.t)
	if r != 0 {
		return timer.GetLoop().LastError().Error()
	}
	return nil
}

func (timer *Timer) Again() (err error) {
	r := C.uv_timer_again(timer.t)
	if r != 0 {
		return timer.GetLoop().LastError().Error()
	}
	return nil
}

func (timer *Timer) SetRepeat(repeat int64) {
	C.uv_timer_set_repeat(timer.t, C.int64_t(repeat))
}

func (timer *Timer) GetRepeat() int64 {
	return int64(C.uv_timer_get_repeat(timer.t))
}
