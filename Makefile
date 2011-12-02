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

OFLAGS+=-luv

include $(GOROOT)/src/Make.pkg
