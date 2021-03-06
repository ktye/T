// T is a text editor.
package main

import (
	"context"
	"flag"
	"image"
	"image/draw"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"time"

	"github.com/eaburns/T/ui"
	"golang.org/x/exp/shiny/driver/gldriver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/math/f64"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

const tickRate = 20 * time.Millisecond

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func main() {
	gldriver.Main(func(scr screen.Screen) {
		flag.Parse()
		if *cpuprofile != "" {
			f, err := os.Create(*cpuprofile)
			if err != nil {
				log.Fatal("could not create CPU profile: ", err)
			}
			if err := pprof.StartCPUProfile(f); err != nil {
				log.Fatal("could not start CPU profile: ", err)
			}
			defer pprof.StopCPUProfile()
		}
		<-newWindow(context.Background(), scr).done
	})
}

type win struct {
	ctx    context.Context
	cancel func()
	done   chan struct{}

	dpi  float32
	size image.Point
	screen.Window

	win *ui.Win
}

func newWindow(ctx context.Context, scr screen.Screen) *win {
	window, err := scr.NewWindow(nil)
	if err != nil {
		panic(err)
	}
	var e size.Event
	for {
		var ok bool
		if e, ok = window.NextEvent().(size.Event); ok {
			break
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	w := &win{
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
		dpi:    float32(e.PixelsPerPt) * 72.0,
		size:   e.Size(),
		Window: window,
	}
	w.win = ui.NewWin(w.dpi)
	w.win.Resize(w.size)

	go tick(w)
	go poll(scr, w)
	return w
}

func (w *win) Release() { w.cancel() }

type done struct{}

func tick(w *win) {
	ticker := time.NewTicker(tickRate)
	for {
		select {
		case <-ticker.C:
			w.Send(time.Now())
		case <-w.ctx.Done():
			ticker.Stop()
			w.Send(done{})
			return
		}
	}
}

func poll(scr screen.Screen, w *win) {
	var mods [4]bool
	dirty := true
	buf, tex := bufTex(scr, w.size)

	for {
		switch e := w.NextEvent().(type) {
		case done:
			buf.Release()
			tex.Release()
			w.Window.Release()
			close(w.done)
			return

		case time.Time:
			if w.win.Tick() {
				w.Send(paint.Event{})
			}

		case lifecycle.Event:
			if e.To == lifecycle.StageDead {
				w.cancel()
				continue
			}
			w.win.Focus(e.To == lifecycle.StageFocused)

		case size.Event:
			if e.Size() == image.ZP {
				w.cancel()
				continue
			}
			w.size = e.Size()
			w.win.Resize(w.size)
			dirty = true
			if b := tex.Bounds(); b.Dx() < w.size.X || b.Dy() < w.size.Y {
				tex.Release()
				buf.Release()
				buf, tex = bufTex(scr, w.size.Mul(2))
			}

		case paint.Event:
			rect := image.Rectangle{Max: w.size}
			img := buf.RGBA().SubImage(rect).(*image.RGBA)
			w.win.Draw(dirty, img)
			dirty = false
			tex.Upload(image.ZP, buf, buf.Bounds())
			w.Draw(f64.Aff3{
				1, 0, 0,
				0, 1, 0,
			}, tex, tex.Bounds(), draw.Src, nil)
			w.Publish()

		case mouse.Event:
			mouseEvent(w, e)

		case key.Event:
			mods = keyEvent(w, mods, e)
		}
	}
}

func mouseEvent(w *win, e mouse.Event) {
	switch pt := image.Pt(int(e.X), int(e.Y)); {
	case e.Button == mouse.ButtonWheelUp:
		w.win.Wheel(pt, 0, 1)

	case e.Button == mouse.ButtonWheelDown:
		w.win.Wheel(pt, 0, -1)

	case e.Button == mouse.ButtonWheelLeft:
		w.win.Wheel(pt, -1, 0)

	case e.Button == mouse.ButtonWheelRight:
		w.win.Wheel(pt, 1, 0)

	case e.Direction == mouse.DirNone:
		w.win.Move(pt)

	case e.Direction == mouse.DirPress:
		w.win.Click(pt, int(e.Button))

	case e.Direction == mouse.DirRelease:
		w.win.Click(pt, -int(e.Button))

	case e.Direction == mouse.DirStep:
		w.win.Click(pt, int(e.Button))
		w.win.Click(pt, -int(e.Button))
	}
}

func keyEvent(w *win, mods [4]bool, e key.Event) [4]bool {
	if e.Direction == key.DirNone {
		e.Direction = key.DirPress
	}
	if e.Direction == key.DirPress && dirKeyCode[e.Code] {
		dirKey(w, e)
		return mods
	}

	switch {
	case e.Code == key.CodeDeleteBackspace:
		e.Rune = '\b'
	case e.Code == key.CodeDeleteForward:
		e.Rune = 0x7f
	case e.Rune == '\r':
		e.Rune = '\n'
	}
	if e.Rune > 0 {
		if e.Direction == key.DirPress {
			w.win.Rune(e.Rune)
		}
		return mods
	}

	return modKey(w, mods, e)
}

var dirKeyCode = map[key.Code]bool{
	key.CodeUpArrow:    true,
	key.CodeDownArrow:  true,
	key.CodeLeftArrow:  true,
	key.CodeRightArrow: true,
	key.CodePageUp:     true,
	key.CodePageDown:   true,
	key.CodeHome:       true,
	key.CodeEnd:        true,
}

func dirKey(w *win, e key.Event) {
	switch e.Code {
	case key.CodeUpArrow:
		w.win.Dir(0, -1)

	case key.CodeDownArrow:
		w.win.Dir(0, 1)

	case key.CodeLeftArrow:
		w.win.Dir(-1, 0)

	case key.CodeRightArrow:
		w.win.Dir(1, 0)

	case key.CodePageUp:
		w.win.Dir(0, -2)

	case key.CodePageDown:
		w.win.Dir(0, 2)

	case key.CodeHome:
		w.win.Dir(0, math.MinInt16)

	case key.CodeEnd:
		w.win.Dir(0, math.MaxInt16)

	default:
		panic("impossible")
	}
}

func modKey(w *win, mods [4]bool, e key.Event) [4]bool {
	var newMods [4]bool
	if e.Modifiers&key.ModShift != 0 {
		newMods[1] = true
	}
	if e.Modifiers&key.ModAlt != 0 {
		newMods[2] = true
	}
	if e.Modifiers&key.ModMeta != 0 ||
		e.Modifiers&key.ModControl != 0 {
		newMods[3] = true
	}
	for i := 0; i < len(newMods); i++ {
		if newMods[i] != mods[i] {
			m := i
			if !newMods[i] {
				m = -m
			}
			w.win.Mod(m)
			mods = newMods
			break
		}
	}
	return mods
}

func bufTex(scr screen.Screen, sz image.Point) (screen.Buffer, screen.Texture) {
	buf, err := scr.NewBuffer(sz)
	if err != nil {
		panic(err)
	}
	tex, err := scr.NewTexture(sz)
	if err != nil {
		panic(err)
	}
	return buf, tex
}
