include $(GOROOT)/src/Make.inc

TARG     = github.com/mattn/go-uv
CGOFILES = \
	callback.go  \
	tcp.go  \
	udp.go  \
	pipe.go  \
	timer.go  \
	loop.go  \
	version.go  \
	error.go  \
	addr.go  \
	idle.go  \
	memory.go \
	process.go \

#CGO_OFILES += $(shell find $(HOME)/dev/libuv/src -name "*.o")

include $(GOROOT)/src/Make.pkg
