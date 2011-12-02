package uv

/*
#include <uv/uv.h>
*/
import "C"

type Error struct {
	e C.uv_err_t
}

func LastError() *Error {
	return &Error{C.uv_last_error(C.uv_default_loop())}
}

func (err *Error) String() string {
	return C.GoString(C.uv_strerror(err.e))
}

func (err *Error) Name() string {
	return C.GoString(C.uv_err_name(err.e))
}
