include $(GOROOT)/src/Make.inc

TARG     = github.com/mattn/go-uv
CGOFILES = uv.go
OFLAGS+=-luv

include $(GOROOT)/src/Make.pkg
