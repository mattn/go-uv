package uv

/*
#include <stdlib.h>
#include <uv/uv.h>

extern void __uv_idle_cb(void* p, int status);
static void _uv_idle_cb(uv_idle_t* idle, int status) {
	__uv_idle_cb(idle->data, status);
}

static int _uv_idle_start(uv_idle_t* idle) {
	return uv_idle_start(idle, _uv_idle_cb);
}

UV_EXTERN int uv_idle_start(uv_idle_t* idle, uv_idle_cb cb);

UV_EXTERN int uv_idle_stop(uv_idle_t* idle);
*/
import "C"
import "errors"
import "unsafe"

type Idle struct {
	i *C.uv_idle_t
}

type idle_callback_info struct {
	idle_cb      func(int)
}

func IdleInit() (idle *Idle, err error) {
	var i C.uv_idle_t

	r := C.uv_idle_init(C.uv_default_loop(), &i)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return nil, errors.New(C.GoString(C.uv_strerror(e)))
	}
	i.data = unsafe.Pointer(&idle_callback_info{})
	return &Idle{&i}, nil
}

func (idle *Idle) Start(cb func(int)) (err error) {
	cbi := (*idle_callback_info)(idle.i.data)
	cbi.idle_cb = cb
	r := C._uv_idle_start(idle.i)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

func (idle *Idle) Stop() (err error) {
	r := C.uv_idle_stop(idle.i)
	if r != 0 {
		e := C.uv_last_error(C.uv_default_loop())
		return errors.New(C.GoString(C.uv_strerror(e)))
	}
	return nil
}

//export __uv_idle_cb
func __uv_idle_cb(p unsafe.Pointer, status int) {
	cbi := (*idle_callback_info)(p)
	if cbi.idle_cb != nil {
		cbi.idle_cb(status)
	}
}
