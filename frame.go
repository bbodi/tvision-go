package tvision

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Frame struct {
	title     string
	HasBorder bool
}

func (self *Frame) drawTitle(view *View, drawBuffer *DrawBuffer) {
	frameW := view.Rect.W
	titleW := len(self.title)
	titleX := frameW/2 - titleW/2
	drawBuffer.DrawWord(titleX, 0, self.title, view.TitleColor(), view.FrameColor())
}

func (self *Frame) drawHorizontalBorder(x, y, w int, drawBuffer *DrawBuffer, view *View) {
	line := strings.Repeat("-", w)
	drawBuffer.DrawWord(x, y, line, view.BorderColor(), view.FrameColor())
}

func (self *Frame) DrawVerticalBorder(x, y, h int, drawBuffer *DrawBuffer, view *View) {
	for i := 0; i < h; i++ {
		drawBuffer.DrawWord(x, y+i, "|", view.BorderColor(), view.FrameColor())
	}
}

func (self *Frame) drawCorners(view *View, drawBuffer *DrawBuffer) {
	fgcolor := view.BorderColor()
	bgcolor := view.FrameColor()
	drawBuffer.DrawWord(0, 0, "+", fgcolor, bgcolor)
	drawBuffer.DrawWord(view.Rect.W-1, 0, "+", fgcolor, bgcolor)
	drawBuffer.DrawWord(view.Rect.W-1, view.Rect.H-1, "+", fgcolor, bgcolor)
	drawBuffer.DrawWord(0, view.Rect.H-1, "+", fgcolor, bgcolor)
}

func (self *Frame) Draw(view *View, drawBuffer *DrawBuffer, bg Color) {
	x, y := 0, 0
	w, h := view.Rect.W, view.Rect.H

	for i := 0; i < h; i++ {
		drawBuffer.SetCellsChBg(x, y+i, w, ' ', bg)
	}
	if self.HasBorder {
		self.drawBorder(view, drawBuffer)
	}
}

func (self *Frame) drawBorder(view *View, drawBuffer *DrawBuffer) {
	self.drawHorizontalBorder(1, 0, view.Rect.W-2, drawBuffer, view)
	self.drawHorizontalBorder(1, view.Rect.H-1, view.Rect.W-2, drawBuffer, view)
	self.DrawVerticalBorder(0, 1, view.Rect.H-2, drawBuffer, view)
	self.DrawVerticalBorder(view.Rect.W-1, 1, view.Rect.H-2, drawBuffer, view)
	self.drawCorners(view, drawBuffer)
	self.drawTitle(view, drawBuffer)
}

func (self *Frame) HandleEvent(event *Event, view *View) {

}

func (self *Frame) calcFrameSize(view *View, items []string) (int, int) {
	w := view.Rect.W
	for _, item := range items {
		if len(item)+2 > w {
			w = len(item) + 2
		}
	}
	if w < utf8.RuneCountInString(self.title)+4 {
		w = utf8.RuneCountInString(self.title) + 4
	}
	return w, len(items) + 1
}

func itemsToString(items []interface{}) []string {
	strings := make([]string, len(items))
	for i, item := range items {
		strings[i] = fmt.Sprintf("%s", item)
	}
	return strings
}
