package tvision

import (
	"fmt"
)

type ComboBox struct {
	originalRect          Rect
	frame                 Frame
	items                 []interface{}
	filteredItems         []interface{}
	allowInputCustomValue bool
	opened                bool
	selectedIndex         int
	input                 string
	selectedItem          interface{}
	Editable              bool
	cursorPos             int
	drawBuffer            *DrawBuffer
	toggleShowCursorTicks int
	showCursor            bool
}

func CreateComboBox(w int) (*View, *ComboBox) {
	view := new(View)
	combo := new(ComboBox)
	combo.frame.HasBorder = false
	combo.drawBuffer = CreateDrawBuffer(w, 1)
	view.Widget = combo
	view.Rect.W = w
	view.Rect.H = 1
	return view, combo
}

func (self *ComboBox) AddItem(item interface{}) {
	self.items = append(self.items, item)
}

func (self *ComboBox) String() string {
	return "ComboBox"
}

func (self *ComboBox) HandleEvent(event *Event, view *View) {
	switch event.Type {
	case EvKey:
		self.handleKey(view, event)
	case EvLostFocus:
		if self.opened {
			self.close(view)
		}
	case EvGetFocus:
		self.toggleShowCursorTicks = 10
	case EvTick:
		if view.focused == false {
			return
		}
		self.toggleShowCursorTicks--
		if self.toggleShowCursorTicks <= 0 {
			self.toggleShowCursorTicks = 10
			self.showCursor = !self.showCursor
			view.Modified()
		}
	}
}

func (self *ComboBox) handleKey(view *View, event *Event) {
	switch event.Key {
	case KeyEsc:
		if self.opened {
			self.close(view)
		}
		event.SetProcessed()
		return
	case KeyArrowDown:
		self.handleDownKey(view)
		event.SetProcessed()
		return
	case KeyArrowUp:
		self.handleUpKey(view)
		event.SetProcessed()
		return
	case KeyEnter:
		self.handleEnter(view)
		event.SetProcessed()
		return
	case KeyTab:
		if self.opened {
			self.handleSelect(view)
			event.SetProcessed()
		}
		return
	}
	if self.Editable == false {
		return
	}
	if self.opened {
		self.close(view)
	}
	ch := event.Ch
	self.input = self.input + string(ch)
	view.Modified()
	event.SetProcessed()
}

func (self *ComboBox) handleEnter(view *View) {
	if self.opened {
		self.handleSelect(view)
	} else {
		if self.Editable && self.input != "" {
			self.AddItem(self.input)
		}
		self.open(view)
	}
}

func (self *ComboBox) handleSelect(view *View) {
	if self.selectedIndex < len(self.items) {
		self.selectedItem = self.items[self.selectedIndex]
		self.input = ""
	}
	self.close(view)
}

func (self *ComboBox) Draw(view *View) *DrawBuffer {
	self.frame.Draw(view, self.drawBuffer, view.ComboFrameColor())
	w, _ := view.Rect.W, view.Rect.H
	if self.opened == false {
		self.drawBuffer.SetCh(w-1, 0, '▼')
		self.drawCurrentItem(view)
		self.drawCursor(view)
	} else {
		self.drawBuffer.SetCh(w-1, 0, '▲')
		self.drawCurrentItem(view)
		self.drawItems(view, itemsToString(self.items), self.drawBuffer)
	}
	return self.drawBuffer
}

func (self *ComboBox) drawCursor(view *View) {
	if self.showCursor == false {
		return
	}
	self.drawBuffer.SetCh(self.cursorPos, 0, '_')
}

func (self *ComboBox) drawCurrentItem(view *View) {
	var drawnText string
	if len(self.input) != 0 {
		drawnText = self.input
	} else if self.selectedItem != nil {
		drawnText = fmt.Sprintf("%s", self.selectedItem)
	}
	if len(drawnText) != 0 {
		w, _ := view.Rect.W, view.Rect.H
		strX := (w - len(drawnText) - 2)
		self.drawBuffer.DrawWordFg(strX, 0, drawnText, view.TextColor())
	}
}

func (self *ComboBox) drawItems(view *View, items []string, drawBuffer *DrawBuffer) {
	for i, item := range items {
		x := 1
		fg := ItemFgColor(self.selectedIndex == i)
		bg := ItemBgColor(self.selectedIndex == i)
		drawBuffer.DrawWord(x, 1+i, item, fg, bg)
	}
}

func (self *ComboBox) handleDownKey(view *View) {
	if self.opened == false {
		self.open(view)
	} else {
		self.moveSelectionDown(view)
	}
}

func (self *ComboBox) handleUpKey(view *View) {
	if self.opened {
		self.moveSelectionUp(view)
	}
}

func (self *ComboBox) moveSelectionDown(view *View) {
	self.selectedIndex++
	if self.selectedIndex >= len(self.items) {
		self.selectedIndex = 0
	}
	view.Modified()
}

func (self *ComboBox) moveSelectionUp(view *View) {
	self.selectedIndex--
	if self.selectedIndex < 0 {
		self.selectedIndex = len(self.items) - 1
		if self.selectedIndex < 0 {
			self.selectedIndex = 0
		}
	}
	view.Modified()
}

func (self *ComboBox) open(view *View) {
	self.originalRect = view.Rect
	self.opened = true
	self.changeViewSizeForItemList(view, self.items)
	if self.selectedIndex >= len(self.items) {
		self.selectedIndex = 0
	}
	view.Modified()
}

func (self *ComboBox) changeViewSizeForItemList(view *View, items []interface{}) int {
	strings := itemsToString(items)
	w, h := self.frame.calcFrameSize(view, strings)
	view.Rect = Rect{view.Rect.X, view.Rect.Y, w, h}
	self.drawBuffer = CreateDrawBuffer(w, h)
	return w
}

func (self *ComboBox) close(view *View) {
	view.Rect = self.originalRect
	self.drawBuffer = CreateDrawBuffer(view.Rect.W, 1)
	self.opened = false
	view.Modified()
}
