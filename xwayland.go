package wlr

// #include <wlr/xwayland.h>
import "C"

import "unsafe"

type XWayland struct {
	p *C.struct_wlr_xwayland
}

type XWaylandSurface struct {
	p *C.struct_wlr_xwayland_surface
}

func NewXWayland(display *Display, compositor *Compositor, lazy bool) *XWayland {
	p := C.wlr_xwayland_create(display.p, compositor.p, C.bool(lazy))
	return &XWayland{p: p}
}

func (x XWayland) Destroy() {
	C.wlr_xwayland_destroy(x.p)
}

func (x *XWayland) OnNewSurface(cb func(*XWaylandSurface)) func() {
	lis := newListener(unsafe.Pointer(x.p), func(lis *wlrlis, data unsafe.Pointer) {
		surface := &XWaylandSurface{p: (*C.struct_wlr_xwayland_surface)(data)}
		trackObject(unsafe.Pointer(surface.p), &surface.p.events.destroy)
		//man.add(unsafe.Pointer(surface.p.surface), &surface.p.surface.events.destroy, func(data unsafe.Pointer) {
		//	man.delete(unsafe.Pointer(surface.p.surface))
		//})
		cb(surface)
	})
	C.wl_signal_add(&x.p.events.new_surface, lis)
	return func() {
		removeListener(lis)
	}
}

func (x *XWayland) SetCursor(img *XCursorImage) {
	C.wlr_xwayland_set_cursor(x.p, img.p.buffer, img.p.width*4, img.p.width, img.p.height, C.int32_t(img.p.hotspot_x), C.int32_t(img.p.hotspot_y))
}

func (s *XWaylandSurface) Surface() *Surface {
	return &Surface{p: s.p.surface}
}

func (s XWaylandSurface) Geometry() Box {
	return Box{
		X:      int(s.p.x),
		Y:      int(s.p.y),
		Width:  int(s.p.width),
		Height: int(s.p.height),
	}
}

func (s XWaylandSurface) Configure(x int16, y int16, width uint16, height uint16) {
	C.wlr_xwayland_surface_configure(s.p, C.int16_t(x), C.int16_t(y), C.uint16_t(width), C.uint16_t(height))
}

func (s *XWaylandSurface) OnMap(cb func(*XWaylandSurface)) func() {
	lis := newListener(unsafe.Pointer(s.p), func(lis *wlrlis, data unsafe.Pointer) {
		cb(s)
	})
	C.wl_signal_add(&s.p.events._map, lis)
	return func() {
		removeListener(lis)
	}
}

func (s *XWaylandSurface) OnUnmap(cb func(*XWaylandSurface)) func() {
	lis := newListener(unsafe.Pointer(s.p), func(lis *wlrlis, data unsafe.Pointer) {
		cb(s)
	})
	C.wl_signal_add(&s.p.events.unmap, lis)
	return func() {
		removeListener(lis)
	}
}

func (s *XWaylandSurface) OnDestroy(cb func(*XWaylandSurface)) func() {
	lis := newListener(unsafe.Pointer(s.p), func(lis *wlrlis, data unsafe.Pointer) {
		cb(s)
	})
	C.wl_signal_add(&s.p.events.destroy, lis)
	return func() {
		removeListener(lis)
	}
}

func (s *XWaylandSurface) OnRequestMove(cb func(surface *XWaylandSurface)) func() {
	lis := newListener(unsafe.Pointer(s.p), func(lis *wlrlis, data unsafe.Pointer) {
		cb(s)
	})
	C.wl_signal_add(&s.p.events.request_move, lis)
	return func() {
		removeListener(lis)
	}
}

func (s *XWaylandSurface) OnRequestResize(cb func(surface *XWaylandSurface, edges Edges)) func() {
	lis := newListener(unsafe.Pointer(s.p), func(lis *wlrlis, data unsafe.Pointer) {
		event := (*C.struct_wlr_xwayland_resize_event)(data)
		cb(s, Edges(event.edges))
	})
	C.wl_signal_add(&s.p.events.request_resize, lis)
	return func() {
		removeListener(lis)
	}
}

func (s *XWaylandSurface) OnRequestConfigure(cb func(surface *XWaylandSurface, x int16, y int16, width uint16, height uint16)) func() {
	lis := newListener(unsafe.Pointer(s.p), func(lis *wlrlis, data unsafe.Pointer) {
		event := (*C.struct_wlr_xwayland_surface_configure_event)(data)
		cb(s, int16(event.x), int16(event.y), uint16(event.width), uint16(event.height))
	})
	C.wl_signal_add(&s.p.events.request_configure, lis)
	return func() {
		removeListener(lis)
	}
}
