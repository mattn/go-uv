package uv

/*
#include <uv.h>
*/
import "C"

type Loop struct {
	l *C.uv_loop_t
}

func Run() {
	DefaultLoop().Run()
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
	C.uv_run(loop.l, C.UV_RUN_DEFAULT)
}

func (loop *Loop) RunOnce() {
	C.uv_run(loop.l, C.UV_RUN_ONCE)
}

func (loop *Loop) UpdateTime() {
	C.uv_update_time(loop.l)
}

func (loop *Loop) Now() int64 {
	return int64(C.uv_now(loop.l))
}
