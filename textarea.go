package tvision

import (
	"strings"
)

type TextRegion struct {
	fgColor Color
	bgColor Color
	text    string
}

func DefaultText(txt string) TextRegion {
	return TextRegion{ColorWhite, ColorBlack, txt}
}

type TextArea struct {
	frame  Frame
	lines  [][]TextRegion
	View   *View
	buffer *DrawBuffer
}

func (self *TextArea) editingRect(view *View) Rect {
	return view.Rect.Grow(-2, -2).Move(1, 1)
}

func CreateTextArea(w, h int, title string) (*View, *TextArea) {
	view := new(View)
	text := new(TextArea)
	text.frame.title = title
	view.Widget = text
	text.View = view
	text.buffer = CreateDrawBuffer(w, h)
	view.Rect.W = w
	view.Rect.H = h
	return view, text
}

func (self *TextArea) drawWord(view *View, x, y int, word string, fg, bg Color, drawBuffer *DrawBuffer) (int, int) {
	tooLong := x+len(word) > view.Rect.W
	if tooLong {
		tooLongForARow := len(word) > view.Rect.W
		if tooLongForARow {
			word = word[:view.Rect.W]
		} else {
			startX := 1
			x = startX
			y += 1
		}
	}
	drawBuffer.DrawWord(x, y, word, fg, bg)
	x += len(word)
	return x, y
}

func (self *TextArea) drawRegion(view *View, x, y int, region TextRegion, drawBuffer *DrawBuffer) (int, int) {
	for _, word := range strings.Fields(region.text) {
		x, y = self.drawWord(view, x, y, word, region.fgColor, region.bgColor, drawBuffer)
	}
	return x, y
}

func (self *TextArea) String() string {
	return "TextArea"
}

func (self *TextArea) Draw(view *View) *DrawBuffer {
	self.frame.Draw(view, self.buffer, view.FrameColor())
	x, y := 1, 1
	for _, line := range self.lines {
		for _, region := range line {
			if y >= view.Rect.H {
				break
			}
			x, y = self.drawRegion(view, x, y, region, self.buffer)
		}
		y += 1
		x = 1
	}
	return self.buffer
}

func (self *TextArea) HandleEvent(event *Event, view *View) {
}

func (self *TextArea) PushBack(region TextRegion) {
	self.lines = append(self.lines, []TextRegion{region})
	self.View.Modified()
}

func (self *TextArea) PushFront(region TextRegion) {
	self.lines = append([][]TextRegion{[]TextRegion{region}}, self.lines...)
	self.View.Modified()
}
