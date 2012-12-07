package uv

/*
#include <uv.h>
*/
import "C"
import "fmt"

func Version() string {
    return fmt.Sprintf("%d.%d", C.UV_VERSION_MAJOR, C.UV_VERSION_MINOR)
}
