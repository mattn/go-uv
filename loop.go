package uv

/*
#include <uv/uv.h>
*/
import "C"

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

func (loop *Loop) RunOnce() {
	C.uv_run_once(loop.l)
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

func (loop *Loop) LastError() *Error {
	return &Error{C.uv_last_error(loop.l)}
}
