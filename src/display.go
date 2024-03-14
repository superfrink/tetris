package main

import (
	"sync"

	"github.com/nsf/termbox-go"
)

// DOC: Display wraps a mutex around termbox writes.
type Display struct {
	mu sync.Mutex
}

func (d *Display) Clear(a termbox.Attribute, b termbox.Attribute) {
	d.mu.Lock()
	termbox.Clear(a, b)
	d.mu.Unlock()
}

func (d *Display) Close() {
	d.mu.Lock()
	termbox.Close()
	d.mu.Unlock()
}

func (d *Display) Flush() {
	d.mu.Lock()
	termbox.Flush()
	d.mu.Unlock()
}

func (d *Display) Init() error {
	d.mu.Lock()
	val := termbox.Init()
	d.mu.Unlock()
	return val
}

func (d *Display) PollEvent() termbox.Event {
	val := termbox.PollEvent()
	return val
}

func (d *Display) TBPrint(y int, x int, str string) {
	d.mu.Lock()
	// tbprint is based on from https://github.com/jjinux/gotetris/
	for _, c := range str {
		termbox.SetCell(x, y, c, termbox.ColorWhite, termbox.ColorBlack)
		x++
	}
	d.mu.Unlock()
}
