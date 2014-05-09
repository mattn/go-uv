package uv

/*
#include <uv.h>
*/
import "C"

type Error struct {
	e int
}

func (err *Error) String() string {
	var error C.uv_err_t
	error.code = C.uv_err_code(err.e)
	return C.GoString(C.uv_strerror(error))
}

func (err *Error) Error() string {
	return err.String()
}

func (err *Error) Name() string {
	var error C.uv_err_t
	error.code = C.uv_err_code(err.e)
	return C.GoString(C.uv_err_name(error))
}
