package wlr

/*
#include <wayland-server-core.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

type Display struct {
	p *C.struct_wl_display
}

func NewDisplay() Display {
	p := C.wl_display_create()
	return Display{p: p}
}

func (d Display) Destroy() {
	C.wl_display_destroy(d.p)
}

func (d Display) OnDestroy(cb func(Display)) Listener {
	lis := newListener(nil, func(lis Listener, data unsafe.Pointer) {
		cb(d)
	})
	C.wl_display_add_destroy_listener(d.p, lis.p)
	return lis
}

func (d Display) Run() {
	C.wl_display_run(d.p)
}

func (d Display) Terminate() {
	C.wl_display_terminate(d.p)
}

func (d Display) EventLoop() EventLoop {
	p := C.wl_display_get_event_loop(d.p)
	evl := EventLoop{p: p}
	return evl
}

func (d Display) AddSocketAuto() (string, error) {
	socket := C.wl_display_add_socket_auto(d.p)
	if socket == nil {
		return "", errors.New("can't auto add wayland socket")
	}

	return C.GoString(socket), nil
}

func (d Display) FlushClients() {
	C.wl_display_flush_clients(d.p)
}
