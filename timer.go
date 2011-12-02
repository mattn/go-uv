package uv

/*
#include <stdlib.h>
#include <uv/uv.h>

extern void __uv_timer_cb(void* p, int status);
static void _uv_timer_cb(uv_timer_t* handle, int status) {
	__uv_timer_cb(handle->data, status);
}

static int _uv_timer_start(uv_timer_t* handle, int64_t timeout, int64_t repeat) {
	return uv_timer_start(handle, _uv_timer_cb, timeout, repeat);
}
*/
import "C"
import "errors"
import "unsafe"

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
