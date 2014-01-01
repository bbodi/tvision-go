package tvision

import (
	"strconv"
)

type Window struct {
	frame          Frame
	drawBuffer     *DrawBuffer
	Resizeable     bool
	Closeable      bool
	HideResizeable bool
	HideCloseable  bool
}

func (self *Window) editingRect(view *View) Rect {
	return view.Rect.Grow(-2, -2).Move(1, 1)
}

func CreateWindow(w, h int, title string) (*View, *Window) {
	view := new(View)
	win := new(Window)
	win.frame.title = title
	win.frame.HasBorder = true
	win.drawBuffer = CreateDrawBuffer(w, h)
	view.Widget = win
	view.Rect.W = w
	view.Rect.H = h
	return view, win
}

func (self *Window) String() string {
	return "Window"
}

func (self *Window) Draw(view *View) *DrawBuffer {
	self.frame.Draw(view, self.drawBuffer, view.FrameColor())
	if self.HideCloseable == false {
		self.drawBuffer.DrawWordFg(view.Rect.W-5, 0, "[", view.BorderColor())
		self.drawBuffer.DrawWordFg(view.Rect.W-4, 0, "X", BorderActionColor(self.Closeable))
		self.drawBuffer.DrawWordFg(view.Rect.W-3, 0, "]", view.BorderColor())
	}
	if self.HideResizeable == false {
		fontSizeStr := strconv.Itoa(int(view.FontSize()))
		lenStr := len(fontSizeStr)
		self.drawBuffer.DrawWordFg(view.Rect.W-10, 0, "[", view.BorderColor())
		self.drawBuffer.DrawWordFg(view.Rect.W-9, 0, fontSizeStr, BorderActionColor(self.Resizeable))
		self.drawBuffer.DrawWordFg(view.Rect.W-9+lenStr, 0, "]", view.BorderColor())
	}
	return self.drawBuffer
}

func (self *Window) HandleEvent(event *Event, view *View) {
	switch event.Type {
	case EvKey:
		switch event.Key {
		case KeyTab:
			view.NextView()
			event.SetProcessed()
			view.Modified()
		}
	case EvMouse:
		self.handleClick(view, event)
	}
}

func (self *Window) handleClick(view *View, event *Event) {
	if clickedInFontChange(view, event.LocalMouseX, event.LocalMouseY) && self.Resizeable {
		event.SetProcessed()
		boxView, box := CreateSelectBox("Fonts")
		box.AddItem(14, CmdOk, 14)
		box.AddItem(16, CmdOk, 16)
		box.AddItem(18, CmdOk, 18)
		box.AddItem(20, CmdOk, 20)
		box.AddItem(22, CmdOk, 22)
		box.AddItem(24, CmdOk, 24)
		result := view.ExecuteView(view.Rect.W-10, 1, boxView)
		if result.Cmd == CmdOk {
			view.SetFontSize(result.Data.(int))
			view.Modified()
		}
	}
}

func clickedInFontChange(view *View, x, y int) bool {
	fontSizeStr := strconv.Itoa(int(view.FontSize()))
	lenStr := len(fontSizeStr)
	return y == 0 && x > view.Rect.W-10 && x < view.Rect.W-9+lenStr
}
