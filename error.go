package uv

/*
#include <uv.h>
*/
import "C"

type Error struct {
	e int
}

func (err *Error) String() string {
	return C.GoString(C.uv_strerror(C.int(err.e)))
}

func (err *Error) Error() string {
	return err.String()
}

func (err *Error) Name() string {
	return C.GoString(C.uv_err_name(C.int(err.e)))
}
