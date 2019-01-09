package ui

import (
	"image/color"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	// framePx is the pixel-width of the lines
	// drawn between columns and rows.
	framePx = 1

	// textPadPx is the pixel-width of the padding
	// between the left and right side of a text box
	// and its text.
	textPadPx = 7

	// cursorWidthPx is the pixel-width of the cursor.
	cursorWidthPx = 4

	// colText is the default column background text.
	colText = "AddCol DelCol Add\n"

	// tagText is the default tag text.
	tagText = " | Del"
)

var (
	// defaultFont is the default font.
	defaultFont, _ = truetype.Parse(goregular.TTF)

	// defaultFontSize is the default font size in points.
	defaultFontSize = 11

	// fg is the text foreground color.
	fg = color.RGBA{R: 0x10, G: 0x28, B: 0x34, A: 0xFF}

	// colBG is the column background color.
	colBG = color.White

	// padColor is the padding color.
	// TODO: remove padColor once it's replaced by frameColor.
	padColor = color.RGBA{R: 0x10, G: 0x28, B: 0x34, A: 0xFF}

	// frameColor is the lines drawn between columns and rows.
	frameColor = fg

	// tagBG is the tag background color.
	tagBG = color.RGBA{R: 0xCF, G: 0xE0, B: 0xF7, A: 0xFF}

	// bodyBG is a body background color.
	bodyBG = color.RGBA{R: 0xFA, G: 0xF0, B: 0xE6, A: 0xFF}

	// hiBG1, hiBG2, and hiBG2 are the background colors
	// of 1-, 2-, and 3-click highlighted text.
	hiBG1 = color.RGBA{R: 0xCC, G: 0xCD, B: 0xAC, A: 0xFF}
	hiBG2 = color.RGBA{R: 0xEC, G: 0x90, B: 0x7F, A: 0xFF}
	hiBG3 = color.RGBA{R: 0xB7, G: 0xE5, B: 0xB2, A: 0xFF}
)
