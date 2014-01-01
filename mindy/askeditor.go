package main

import (
	"github.com/bbodi/tvision"
)

type textRangeType int

type textRange struct {
	from, to int
	typ      textRangeType
}

type AskEditor struct {
	frame            tvision.Frame
	lines            []string
	ranges           []textRangeType
	cursorX, cursorY int
	drawBuffer       *tvision.DrawBuffer
}

func CreateAskEditor(w, h int) (*tvision.View, *AskEditor) {
	view := new(tvision.View)
	editor := new(AskEditor)
	editor.frame.HasBorder = false
	editor.drawBuffer = tvision.CreateDrawBuffer(w, h)
	view.Widget = editor
	view.Rect.W = w
	view.Rect.H = h
	return view, editor
}

func (self *AskEditor) String() string {
	return "AskEditor"
}

func (self *AskEditor) HandleEvent(event *tvision.Event, view *tvision.View) {
	if event.Type != tvision.EvKey {
		return
	}
	switch event.Key {
	case tvision.KeyEsc:
		event.SetProcessed()
		return
	case tvision.KeyArrowDown:
		self.handleDownKey(view)
		event.SetProcessed()
		return
	case tvision.KeyArrowUp:
		self.handleUpKey(view)
		event.SetProcessed()
		return
	case tvision.KeyEnter:
		self.handleEnter(view)
		event.SetProcessed()
		return
	case tvision.KeyTab:
		//tvision.ClearEvent(event)
		return
	}
	//ch := event.Event.Ch
	//line := self.lines[self.cursorY]
	//self.lines[self.cursorY][self.cursorX] = ch
	event.SetProcessed()
	view.Modified()
}

func (self *AskEditor) handleEnter(view *tvision.View) {

}

func (self *AskEditor) Draw(view *tvision.View) *tvision.DrawBuffer {
	self.frame.Draw(view, self.drawBuffer, view.ComboFrameColor())
	//_, _ := view.Rect.W, view.Rect.H
	for y, line := range self.lines {
		for x, ch := range line {
			self.drawBuffer.SetCh(x, y, ch)
		}
	}
	return self.drawBuffer
}

func (self *AskEditor) handleDownKey(view *tvision.View) {
	self.cursorY++
	if self.cursorY >= view.Rect.H {
		self.cursorY = view.Rect.H - 1
	}
}

func (self *AskEditor) handleUpKey(view *tvision.View) {
	self.cursorY--
	if self.cursorY < 0 {
		self.cursorY = 0
	}
}
