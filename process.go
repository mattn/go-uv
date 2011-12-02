package uv

/*
#include <uv/uv.h>

static void _uv_exit_cb(uv_process_t*, int exit_status, int term_signal) {
}
*/
import "C"

type ProcessOptions struct {
  Exit_cb func(int, int)
  File string
  Args []string
  Env []string
  Cwd string
  int WindowsVerbatimArguments
  StdinStream *Pipe
  StdoutStream *Pipe
  StderrStream *Pipe
}
