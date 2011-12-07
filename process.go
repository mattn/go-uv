package uv

/*
#include <uv/uv.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "runtime"

type ProcessOptions struct {
  Exit_cb func(*Handle, int, int)
  File string
  Args []string
  Env []string
  Cwd string
  WindowsVerbatimArguments int
  StdinStream *Pipe
  StdoutStream *Pipe
  StderrStream *Pipe
}

func Spawn(loop *Loop, options ProcessOptions) (err error) {
	if loop == nil {
		loop = DefaultLoop()
	}

	var opt C.uv_process_options_t
	defer func() {
		/*
		C.free(unsafe.Pointer(opt.file))
		if len(options.Args) > 0 {
			for n := 0; n < len(options.Args); n++ {
				C.free(unsafe.Pointer(((*[1<<24]*C.char)(unsafe.Pointer(&opt.args)))[n]))
			}
			C.free(unsafe.Pointer(opt.args))
		}
		if len(options.Env) > 0 {
			for n := 0; n < len(options.Env); n++ {
				C.free(unsafe.Pointer(((*[1<<24]*C.char)(unsafe.Pointer(&opt.env)))[n]))
			}
			C.free(unsafe.Pointer(opt.env))
		}
		C.free(unsafe.Pointer(opt.cwd))
		*/
	}()
	if len(options.File) > 0 {
		opt.file = C.CString(options.File)
	}
	if len(options.Args) > 0 {
		opt.args = (**C.char)(C.malloc(C.size_t(4 * (len(options.Args)+1))))
		for n := 0; n < len(options.Args); n++ {
			((*[1<<24]*C.char)(unsafe.Pointer(&opt.args)))[n] = C.CString(options.Args[n])
		}
		((*[1<<24]*C.char)(unsafe.Pointer(&opt.args)))[len(options.Args)] = nil
	}
	if len(options.Env) > 0 {
		opt.env = (**C.char)(C.malloc(C.size_t(4 * (len(options.Env)+1))))
		for n := 0; n < len(options.Args); n++ {
			((*[1<<24]*C.char)(unsafe.Pointer(&opt.env)))[n] = C.CString(options.Args[n])
		}
		((*[1<<24]*C.char)(unsafe.Pointer(&opt.env)))[len(options.Args)] = nil
	}
	if len(options.Cwd) > 0 {
		opt.cwd = C.CString(options.Cwd)
	}
	opt.windows_verbatim_arguments = C.int(options.WindowsVerbatimArguments)
	if options.StdinStream != nil {
		opt.stdin_stream = options.StdinStream.p
	}
	if options.StdoutStream != nil {
		opt.stdout_stream = options.StdoutStream.p
	}
	if options.StderrStream != nil {
		opt.stderr_stream = options.StderrStream.p
	}
	if runtime.GOOS == "windows" {
		opt.windows_verbatim_arguments = 1
	}

	var p C.uv_process_t
	p.data = unsafe.Pointer(&callback_info{exit_cb: options.Exit_cb})
	r := uv_spawn(loop.l, &p, opt)
	if r != 0 {
		return loop.LastError().Error()
	}
	return nil
}
