package tvision

import (
	"fmt"
	"log"
)

type SelectBox struct {
	frame         Frame
	selectedIndex int
	items         []*SelectItem
	view          *View
}

type SelectItem struct {
	name       string
	returnData interface{}
	cmd        Cmd
	childBox   *SelectBox
	childView  *View
}

func (self *SelectBox) String() string {
	return "SelectBox"
}

func CreateSelectBox(title string) (*View, *SelectBox) {
	view := new(View)
	box := new(SelectBox)
	view.Widget = box
	box.view = view
	box.frame.HasBorder = true
	box.frame.title = title
	return view, box
}

func (self *SelectBox) Draw(view *View) *DrawBuffer {
	drawBuffer := CreateDrawBuffer(80, 40)
	self.frame.Draw(view, drawBuffer, view.FrameColor())
	for i, item := range self.items {
		fg := ItemFgColor(i == self.selectedIndex)
		bg := ItemBgColor(i == self.selectedIndex)
		drawBuffer.DrawWord(1, 2+i, item.name, fg, bg)
	}
	return drawBuffer
}

func (self *SelectBox) HandleEvent(event *Event, view *View) {
	log.Print("SelectBox.HandleEvent: ", event)
	switch event.Type {
	case EvKey:
		switch event.Key {
		case KeyArrowDown:
			self.selectedIndex++
			if self.selectedIndex >= len(self.items) {
				self.selectedIndex = 0
			}
			view.Modified()
			ClearEvent(event)
		case KeyArrowUp:
			self.selectedIndex--
			if self.selectedIndex < 0 {
				self.selectedIndex = len(self.items) - 1
			}
			view.Modified()
			ClearEvent(event)
		case KeyEnter:
			self.handleSelectItem(view)
			ClearEvent(event)
		case KeyEsc:
			if view.isExecuting() {
				view.StopExecuting(ExecutingResult{CmdCancel, nil})
				ClearEvent(event)
			}
		}
	case EvMouse:
		self.handleClick(event, view)
	}
}

func (self *SelectBox) handleClick(event *Event, view *View) {
	if view.clickedInMe(event) == false {
		return
	}
	my := event.LocalMouseY
	index := my - 2
	if index < 0 || index >= len(self.items) {
		return
	}
	self.selectedIndex = my - 2
	ClearEvent(event)
	if event.DoubleClick {
		self.handleSelectItem(view)
	} else {
		view.Modified()
	}
}

func (self *SelectBox) handleSelectItem(view *View) {
	selectedItem := self.items[self.selectedIndex]
	if selectedItem.childBox == nil {
		if view.executing {
			view.StopExecuting(ExecutingResult{selectedItem.cmd, selectedItem.returnData})
		} else if selectedItem.cmd != CmdOk {
			view.BroadcastCommand(selectedItem.cmd, selectedItem.returnData)
		}
	} else {
		result := view.ExecuteView(3, 3, selectedItem.childView)
		if result.Cmd == CmdCancel {
			view.Modified()
			return
		}
		if view.executing {
			view.StopExecuting(ExecutingResult{result.Cmd, result.Data})
		} else if result.Cmd != CmdOk {
			view.BroadcastCommand(result.Cmd, result.Data)
		}
	}
}

func (self *SelectBox) stringItemList() []string {
	strings := make([]string, len(self.items))
	for _, item := range self.items {
		strings = append(strings, item.name)
	}
	return strings
}

func (self *SelectBox) AddItem(item interface{}, cmd Cmd, returnValue interface{}) *SelectItem {
	selectItem := new(SelectItem)
	selectItem.name = fmt.Sprintf("%v", item)
	selectItem.cmd = cmd
	selectItem.returnData = returnValue
	self.items = append(self.items, selectItem)
	self.changeViewSizeForItemList(self.view)
	return selectItem
}

func (self *SelectBox) AddEventItem(item interface{}, cmd Cmd) *SelectItem {
	return self.AddItem(item, cmd, nil)
}

func (self *SelectItem) AddSubItem(item interface{}, cmd Cmd, returnValue interface{}) *SelectItem {
	if self.childBox == nil {
		self.childView, self.childBox = CreateSelectBox(self.name)
	}
	return self.childBox.AddItem(item, cmd, returnValue)
}

func (self *SelectBox) changeViewSizeForItemList(view *View) int {
	strings := self.stringItemList()
	w, h := self.frame.calcFrameSize(view, strings)
	view.Rect = Rect{view.Rect.X, view.Rect.Y, w, h}
	return w
}
