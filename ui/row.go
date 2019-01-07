package ui

import (
	"image"
	"image/draw"
)

// The Row interface is implemented by UI elements
// that sit in a column, draw, and react to user input events.
//
// All corrdinates below are relative to the row,
// with 0,0 in the upper left.
type Row interface {
	// Draw draws the element to the image.
	// If dirty is true the element should redraw itself in its entirity.
	// If dirty is false, the element need only redraw
	// parts that have changed since the last call to Draw.
	Draw(dirty bool, img draw.Image)

	// Focus handles a focus state change.
	// The focus is either true (in focus) or false (out of focus).
	Focus(focus bool)

	// Resize handles a resize event.
	Resize(size image.Point)

	// Tick handles periodic ticks
	// sent to the element at regular intervals.
	// This intended to drive asynchronous events
	// in a synchronous manner.
	Tick() bool

	// Move handles mouse cursor moving events.
	// It returns whether the element needs to be redrawn.
	Move(pt image.Point) bool

	// Click handles mouse button events.
	// It returns whether the element needs to be redrawn.
	//
	// The absolute value of the argument indicates the mouse button.
	// A positive value indicates the button was pressed.
	// A negative value indicates the button was released.
	Click(pt image.Point, button int) ([2]int64, bool)

	// Wheel handles mouse wheel events.
	// It returns whether the element needs to be redrawn.
	// 	-y is roll up.
	// 	+y is roll down.
	// 	-x is roll left.
	// 	+x is roll right.
	Wheel(x, y int) bool

	// Dir handles keyboard directional events.
	// It returns whether the element needs to be redrawn.
	//
	// These events are generated by the arrow keys,
	// page up and down keys, and the home and end keys.
	// Exactly one of x or y must be non-zero.
	//
	// If the absolute value is 1, then it is treated as an arrow key
	// in the corresponding direction (x-horizontal, y-vertical,
	// negative-left/up, positive-right/down).
	// If the absolute value is math.MinInt16, it is treated as a home event.
	// If the absolute value is math.MathInt16, it is end.
	// Otherwise, if the value for y is non-zero it is page up/down.
	// Other non-zero values for x are currently ignored.
	//
	// Dir only handles key press events, not key releases.
	Dir(x, y int) bool

	// Mod handles modifier key state change events.
	// It returns whether the element needs to be redrawn.
	//
	// The absolute value of the argument indicates the modifier key.
	// A positive value indicates the key was pressed.
	// A negative value indicates the key was released.
	Mod(m int) bool

	// Rune handles typing events.
	// It returns whether the element needs to be redrawn.
	//
	// The argument is a rune indicating the glyph typed
	// after interpretation by any system-dependent
	// keyboard/layout mapping.
	// For example, if the 'a' key is pressed
	// while the shift key is held,
	// the argument would be the letter 'A'.
	//
	// If the rune is positive, the event is a key press,
	// if negative, a key release.
	Rune(r rune) bool
}