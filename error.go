package uv

/*
#include <uv/uv.h>
*/
import "C"
import "errors"

type Error struct {
	e C.uv_err_t
}

func (err *Error) String() string {
	return C.GoString(C.uv_strerror(err.e))
}

func (err *Error) Error() error {
	return errors.New(err.String())
}

func (err *Error) Name() string {
	return C.GoString(C.uv_err_name(err.e))
}
