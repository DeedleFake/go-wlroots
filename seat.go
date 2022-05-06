package wlr

// #include <stdlib.h>
// #include <wlr/types/wlr_seat.h>
import "C"

import "unsafe"

type Seat struct {
	p *C.struct_wlr_seat
}

type SeatClient struct {
	p *C.struct_wlr_seat_client
}

type SeatKeyboardState struct {
	s C.struct_wlr_seat_keyboard_state
}

type SeatPointerState struct {
	s C.struct_wlr_seat_pointer_state
}

type SeatCapability uint32

const (
	SeatCapabilityPointer  SeatCapability = C.WL_SEAT_CAPABILITY_POINTER
	SeatCapabilityKeyboard SeatCapability = C.WL_SEAT_CAPABILITY_KEYBOARD
	SeatCapabilityTouch    SeatCapability = C.WL_SEAT_CAPABILITY_TOUCH
)

func NewSeat(display Display, name string) Seat {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	p := C.wlr_seat_create(display.p, s)
	return Seat{p: p}
}

func (s Seat) Destroy() {
	C.wlr_seat_destroy(s.p)
}

func (s Seat) OnDestroy(cb func(Seat)) Listener {
	return newListener(&s.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s Seat) OnRequestSetCursor(cb func(client SeatClient, surface Surface, serial uint32, hotspotX int32, hotspotY int32)) Listener {
	return newListener(&s.p.events.request_set_cursor, func(lis Listener, data unsafe.Pointer) {
		event := (*C.struct_wlr_seat_pointer_request_set_cursor_event)(data)
		client := SeatClient{p: event.seat_client}
		surface := Surface{p: event.surface}
		cb(client, surface, uint32(event.serial), int32(event.hotspot_x), int32(event.hotspot_y))
	})
}

func (s Seat) SetCapabilities(caps SeatCapability) {
	C.wlr_seat_set_capabilities(s.p, C.uint32_t(caps))
}

func (s Seat) SetKeyboard(dev InputDevice) {
	C.wlr_seat_set_keyboard(s.p, dev.p)
}

func (s Seat) NotifyPointerButton(time uint32, button uint32, state ButtonState) {
	C.wlr_seat_pointer_notify_button(s.p, C.uint32_t(time), C.uint32_t(button), uint32(state))
}

func (s Seat) NotifyPointerAxis(time uint32, orientation AxisOrientation, delta float64, deltaDiscrete int32, source AxisSource) {
	C.wlr_seat_pointer_notify_axis(s.p, C.uint32_t(time), C.enum_wlr_axis_orientation(orientation), C.double(delta), C.int32_t(deltaDiscrete), C.enum_wlr_axis_source(source))
}

func (s Seat) NotifyPointerEnter(surface Surface, sx float64, sy float64) {
	C.wlr_seat_pointer_notify_enter(s.p, surface.p, C.double(sx), C.double(sy))
}

func (s Seat) NotifyPointerMotion(time uint32, sx float64, sy float64) {
	C.wlr_seat_pointer_notify_motion(s.p, C.uint32_t(time), C.double(sx), C.double(sy))
}

func (s Seat) NotifyPointerFrame() {
	C.wlr_seat_pointer_notify_frame(s.p)
}

func (s Seat) NotifyKeyboardEnter(surface Surface, k Keyboard) {
	C.wlr_seat_keyboard_notify_enter(s.p, surface.p, &k.p.keycodes[0], k.p.num_keycodes, &k.p.modifiers)
}

func (s Seat) NotifyKeyboardModifiers(k Keyboard) {
	C.wlr_seat_keyboard_notify_modifiers(s.p, &k.p.modifiers)
}

func (s Seat) NotifyKeyboardKey(time uint32, keyCode uint32, state KeyState) {
	C.wlr_seat_keyboard_notify_key(s.p, C.uint32_t(time), C.uint32_t(keyCode), C.uint32_t(state))
}

func (s Seat) ClearPointerFocus() {
	C.wlr_seat_pointer_clear_focus(s.p)
}

func (s Seat) Keyboard() Keyboard {
	p := C.wlr_seat_get_keyboard(s.p)
	return Keyboard{p: p}
}

func (s Seat) KeyboardState() SeatKeyboardState {
	return SeatKeyboardState{s: s.p.keyboard_state}
}

func (s Seat) PointerState() SeatPointerState {
	return SeatPointerState{s: s.p.pointer_state}
}

func (s SeatKeyboardState) FocusedSurface() Surface {
	return Surface{p: s.s.focused_surface}
}

func (s SeatPointerState) FocusedSurface() Surface {
	return Surface{p: s.s.focused_surface}
}

func (s SeatPointerState) FocusedClient() SeatClient {
	return SeatClient{p: s.s.focused_client}
}
