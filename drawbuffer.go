package tvision

import (
	"fmt"
)

type Cell struct {
	Ch rune
	Fg Color
	Bg Color
}

type DrawBuffer struct {
	cells []Cell
	w, h  int
}

func (self *DrawBuffer) Cell(x, y int) *Cell {
	i := self.w*y + x
	return &self.cells[i]
}

func (self *DrawBuffer) SetCell(x, y int, ch rune, fg, bg Color) {
	if self.outOfRange(x, y) {
		return
	}
	i := self.w*y + x
	self.cells[i] = Cell{ch, fg, bg}
}

func (self *DrawBuffer) SetCellBg(x, y int, bg Color) {
	if self.outOfRange(x, y) {
		return
	}
	i := self.w*y + x
	self.cells[i].Bg = bg
}

func (self *DrawBuffer) SetCellChBg(x, y int, ch rune, bg Color) {
	if self.outOfRange(x, y) {
		return
	}
	i := self.w*y + x
	self.cells[i].Bg = bg
	self.cells[i].Ch = ch
}

func (self *DrawBuffer) DrawWordFg(x, y int, obj interface{}, fg Color) {
	cx, cy := x, y
	chars := fmt.Sprintf("%s", obj)
	for _, ch := range chars {
		if cx >= self.w {
			cy += 1
		}
		if cy >= self.h {
			cy += 1
			cx = x
		}
		if self.outOfRange(cx, cy) {
			continue
		}
		i := self.w*cy + cx
		self.cells[i].Ch = ch
		self.cells[i].Fg = fg
		cx += 1
	}

}

func (self *DrawBuffer) DrawWordFgBg(x, y int, obj interface{}, fg, bg Color) {
	self.DrawWord(x, y, obj, fg, bg)
}

func (self *DrawBuffer) SetCellsBg(x, y, w int, bg Color) {
	for cx := x; cx < x+w; cx++ {
		self.SetCellBg(cx, y, bg)
	}
}

func (self *DrawBuffer) SetCellsChBg(x, y, w int, ch rune, bg Color) {
	for cx := x; cx < x+w; cx++ {
		self.SetCellChBg(cx, y, ch, bg)
	}
}

func (self *DrawBuffer) DrawWord(x, y int, obj interface{}, fg, bg Color) {
	cx, cy := x, y
	chars := fmt.Sprintf("%s", obj)
	for _, ch := range chars {
		if cx >= self.w {
			cy += 1
		}
		if cy >= self.h {
			cy += 1
			cx = x
		}
		if self.outOfRange(cx, cy) {
			continue
		}
		i := self.w*cy + cx
		self.cells[i] = Cell{ch, fg, bg}
		cx += 1
	}
}

func (self *DrawBuffer) SetCh(x, y int, ch rune) {
	if self.outOfRange(x, y) {
		return
	}
	i := self.w*y + x
	self.cells[i].Ch = ch

}

func (self *DrawBuffer) SetFg(x, y int, fg Color) {
	if self.outOfRange(x, y) {
		return
	}
	i := self.w*y + x
	self.cells[i].Fg = fg
}

func (self *DrawBuffer) SetBg(x, y int, bg Color) {
	if self.outOfRange(x, y) {
		return
	}
	i := self.w*y + x
	self.cells[i].Bg = bg

}

func (self *DrawBuffer) outOfRange(x, y int) bool {
	return x >= self.w || y >= self.h || y < 0 || x < 0
}

func (self *DrawBuffer) Clear(fg, bg Color) {
	for i := range self.cells {
		cell := &self.cells[i]
		cell.Fg = fg
		cell.Bg = bg
		cell.Ch = 0
	}
}

func (self *DrawBuffer) DrawToBackBuffer(view *View) {
	engine.DrawToBackBuffer(self, view)
}

func CreateDrawBuffer(w, h int) *DrawBuffer {
	buffer := new(DrawBuffer)
	buffer.w = w
	buffer.h = h
	buffer.cells = make([]Cell, w*h)
	return buffer
}

type screenRegion struct {
	drawBuffer *DrawBuffer
	view       *View
}

func DrawToBackBuffer(regions []*screenRegion) {
	for _, region := range regions {
		region.drawBuffer.DrawToBackBuffer(region.view)
	}
}
