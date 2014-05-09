package uv

/*
#include <uv.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "errors"

type SockaddrIn interface {
	Name() (string, error)
}

type Sockaddr struct {
	sa C.struct_sockaddr
}

type SockaddrIn4 struct {
	sa C.struct_sockaddr_in
}

type SockaddrIn6 struct {
	sa C.struct_sockaddr_in6
}

func Ip4Addr(host string, port uint16) (SockaddrIn, error) {
	phost := C.CString(host)
	defer C.free(unsafe.Pointer(phost))
	var addr C.struct_sockaddr_in
	addr = C.uv_ip4_addr(phost, C.int(port))
	return &SockaddrIn4 {addr}, nil
}

func Ip6Addr(host string, port uint16) (SockaddrIn, error) {
	phost := C.CString(host)
	defer C.free(unsafe.Pointer(phost))
	var addr C.struct_sockaddr_in6
	addr = C.uv_ip6_addr(phost, C.int(port))
	return &SockaddrIn6 {addr}, nil
}

func (sa *SockaddrIn4) Name() (name string, err error) {
	b := make([]byte, 256)
	r := C.uv_ip4_name(&sa.sa, (*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)));
	if r != 0 {
		var error C.uv_err_t
		error.code = C.uv_err_code(r)
		return "", errors.New(C.GoString(C.uv_strerror(error)))
	}
	return string(b), nil
}

func (sa *SockaddrIn6) Name() (name string, err error) {
	b := make([]byte, 256)
	r := C.uv_ip6_name(&sa.sa, (*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)));
	if r != 0 {
		var error C.uv_err_t
		error.code = C.uv_err_code(r)
		return "", errors.New(C.GoString(C.uv_strerror(error)))
	}
	return string(b), nil
}
