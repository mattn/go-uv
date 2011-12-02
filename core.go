package uv

/*
#include <uv/uv.h>

#cgo darwin LDFLAGS: -luv
#cgo linux LDFLAGS: -luv
#cgo windows LDFLAGS: -luv.dll -lws2_32
*/
import "C"
import "fmt"

func Version() string {
    return fmt.Sprintf("%d.%d", C.UV_VERSION_MAJOR, C.UV_VERSION_MINOR)
}
