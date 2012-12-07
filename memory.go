package uv

/*
#include <uv.h>
*/
import "C"

func GetFreeMemory() uint64 {
	return uint64(C.uv_get_free_memory())
}

func GetTotalMemory() uint64 {
	return uint64(C.uv_get_total_memory())
}
